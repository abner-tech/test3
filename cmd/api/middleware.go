package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/abner-tech/Test3-Api.git/internal/data"
	"github.com/abner-tech/Test3-Api.git/internal/validator"
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

// check if the user is authenticated NOTE: not anonymous
func (a *applicationDependences) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := a.contextGetUser(r)

		if user.IsAnonymous() {
			a.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// check if user is activated
func (a *applicationDependences) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := a.contextGetUser(r)

		if !user.Activated {
			a.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})

	//We pass the activation check middleware to the authentication
	// middleware to call (next) if the authentication check succeeds
	// In other words, only check if the user is activated if they are
	// actually authenticated.
	return a.requireAuthenticatedUser(fn)
}

// check if the user has teh right permissions, we send permissions which is expected as an argument
func (a *applicationDependences) requirePermission(permissionCode string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := a.contextGetUser(r)
		//get all permissions accociated with the usr
		permissions, err := a.permisionsModel.GetAllForUser(user.ID)
		if err != nil {
			a.serverErrorResponse(w, r, err)
			return
		}
		if !permissions.Include(permissionCode) {
			a.notFoundResponse(w, r)
			return
		}
		//everything good
		next.ServeHTTP(w, r)
	}
	return a.requireActivatedUser(fn)
}

func (a *applicationDependences) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		/*
		 this header must be added to the response object or we are defeating the purpose of CORS, Why? browsers
		  want to be fast, so they cache information, If on one resposne we say that appletree.com is a trusted origin,
		  the browser is tempted to cache this, so if later a response comes in from a different origin (evel.com), the
		  browser will be tempted to look in its cache and do what it did for the last response that came in allow it -
		  which is bad and send the same response. such as displaying personal account information. SO, we tell the browser
		  that trusted origins can change for it not to rely on cache memory.

		*/
		w.Header().Add("Vary", "Origin")
		//we check the request origin and verify if its part of the allowed list
		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range a.config.cors.trustedOrigins {
				if origin == a.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
