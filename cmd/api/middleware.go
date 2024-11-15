package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/abner-tech/Comments-Api.git/internal/data"
	"github.com/abner-tech/Comments-Api.git/internal/validator"
	"golang.org/x/time/rate"
)

func (a *applicationDependences) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//defer will be called when the stack unwinds
		defer func() {
			//recover from panic
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "Close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *applicationDependences) rateLimiting(next http.Handler) http.Handler {
	// rate limit struct
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time //remove map entries that are stable
	}

	var mu sync.Mutex
	var clients = make(map[string]*client)
	//a gorutine to remove stale entries from the map
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock() //begin clean-up
			//delete any entry not seen in 3 minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock() //finish cleanup
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get the ip address
		if a.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				a.serverErrorResponse(w, r, err)
				return
			}
			mu.Lock() //exclusive access to the map
			//check if ip address already is in map, if not add it
			_, found := clients[ip]
			if !found {
				clients[ip] = &client{limiter: rate.NewLimiter(
					rate.Limit(a.config.limiter.rps),
					a.config.limiter.burst,
				)}
			}

			//update the last seen of the clients
			clients[ip].lastSeen = time.Now()

			//check the rate limit status
			if !clients[ip].limiter.Allow() {
				mu.Unlock() //no longer needs exclusive access to the map
				a.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock() //others are free to get exclusive access to the map
		}
		next.ServeHTTP(w, r)

	})
}

func (a *applicationDependences) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*This header tells the servers not to cache the response when
		the Authorization header changes. This also means that the server is not
		supposed to serve the same cached data to all users regardless of their
		Authorization values. Each unique user gets their own cache entry*/
		w.Header().Add("Vary", "Authorization")

		/*Get the Authorization Header from the request. It should have the Bearer token*/
		authorizationHeader := r.Header.Get("Authorization")

		//if no authorization header found, then its an anonymous user
		if authorizationHeader == "" {
			r = a.contextSetUser(r, data.AnonymouseUser)
			next.ServeHTTP(w, r)
			return
		}
		/* Bearer token present so parse it. The Bearer token is in the form
		Authorization: Bearer IEYZQUBEMPPAKPOAWTPV6YJ6RM
		We will implement invalidAuthenticationTokenResponse() later */

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}
		//get the actual token
		token := headerParts[1]
		//validatte
		v := validator.New()

		data.ValidatetokenPlaintext(v, token)
		if !v.IsEmpty() {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		//get the user info relatedw with this authentication token
		user, err := a.userModel.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				a.invalidAuthenticationTokenResponse(w, r)
			default:
				a.serverErrorResponse(w, r, err)
			}
			return
		}
		//add the retrieved user info to the context
		r = a.contextSetUser(r, user)
		//call the next handler in the chair
		next.ServeHTTP(w, r)
	})
}
