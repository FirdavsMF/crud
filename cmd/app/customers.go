package app

import (
	"encoding/json"
	"net/http"

	"github.com/FirdavsMF/crud/pkg/customers"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleCustomerRegistration(w http.ResponseWriter, r *http.Request) {

	//we declare the structure of the client for the request
	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//call the function for an error response
		errorWriter(w, http.StatusBadRequest, err)
		return
	}


//Generating a bcrypt hash from a real password
	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		//call the function for an error response
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//and supply the hash in the password field
	item.Password = string(hashed)

	//we synchronize or update the client
	customer, err := s.customerSvc.Save(r.Context(), item)

	//if we receive an error, we reply with an error
	if err != nil {
		//call the function for an error response
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//we call the function for the response in JSON format
	respondJSON(w, customer)
}

func (s *Server) handleCustomerGetToken(w http.ResponseWriter, r *http.Request) {
	//declaring the structure for the request
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	//extracting data from the request
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//call the function for an error response
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	//call the AuthenticateCustomer method from the securitySvc service
	token, err := s.customerSvc.Token(r.Context(), item.Login, item.Password)

	if err != nil {
		//call the function for an error response
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//we call the function for the response in JSON format
	respondJSON(w, map[string]interface{}{"status": "ok", "token": token})

}

func (s *Server) handleCustomerGetProducts(w http.ResponseWriter, r *http.Request) {

	items, err := s.customerSvc.Products(r.Context())
	if err != nil {
		//call the function for an error response
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	respondJSON(w, items)

}
