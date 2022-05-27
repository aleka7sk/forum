package http

import (
	"forum/internal/auth"
	"html/template"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase auth.UseCase
	log     *logrus.Logger
}

type signInResponse struct {
	Token string `json:"token"`
}

func NewHandler(usecase auth.UseCase) *Handler {
	logger := logrus.New()
	return &Handler{
		usecase: usecase,
		log:     logger,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Write([]byte("404 Error"))
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parse main page index: %v", tmpl)
	}
	tmpl.Execute(w, nil)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signin.html")
	if err != nil {
		log.Printf("Error parse main page signin: %v", tmpl)
	}
	if r.Method == "GET" {
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		token, err := h.usecase.SignIn(r.Context(), username, password)
		if err != nil {
			if err == auth.ErrUserNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Неавторизован"))
				h.log.Info("Не авторизован")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Info("Ошибка сервера")
			return
		}
		h.log.Info("Пользователь авторизован")
		http.SetCookie(w, &http.Cookie{Name: "token", Value: token})
		tmpl.ExecuteTemplate(w, "signin.html", signInResponse{Token: token})
	} else {
		h.log.Info("Плохой запрос")
		w.Write([]byte("Плохой запрос"))
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		log.Printf("Error parse main page signup: %v", tmpl)
	}
	if r.Method == "GET" {
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		err := h.usecase.SignUp(r.Context(), username, password)
		if err != nil {
			tmpl.ExecuteTemplate(w, "signup.html", "Ошибка регистрации")
		}
		tmpl.ExecuteTemplate(w, "signup.html", "Регистрация прошла успешно")
	} else {
		w.Write([]byte("Плохой запрос"))
	}
}

func (h *Handler) Private(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.usecase.ParseToken(r.Context(), c.Value)
	if user != nil {
		log.Printf("Parse token error: %v", err)
	}
	tmpl, err := template.ParseFiles("templates/private.html")
	if err != nil {
		log.Printf("Parse template private error: %v", err)
	}

	tmpl.Execute(w, user)
}
