package http

import (
	"forum/internal/auth"
	"forum/models"
	"html/template"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase auth.UseCase
	log     *logrus.Logger
}

type SignInResponse struct {
	Token string `json:"token"`
}

type IsAuth struct {
	IsAuth bool
	// Data   interface{}
	User *models.User
}

func NewHandler(usecase auth.UseCase) *Handler {
	logger := logrus.New()
	return &Handler{
		usecase: usecase,
		log:     logger,
	}
}

func RenderTemplate(tmpl string) []string {
	files := []string{}
	files = append(files, "templates/"+tmpl)
	static_tmpl := []string{"templates/layout/layout.html"}
	files = append(files, static_tmpl...)
	return files
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(RenderTemplate("signin.html")...)
	if err != nil {
		log.Printf("Error parse main page signin: %v", tmpl)
	}
	right := r.Context().Value("rights")
	if right == "auth" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		// tmpl.ExecuteTemplate(w, "signin.html", signInResponse{Token: token})
	} else {
		h.log.Info("Плохой запрос")
		w.Write([]byte("Плохой запрос"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(RenderTemplate("signup.html")...)
	if err != nil {
		log.Printf("Error parse main page signup: %v", tmpl)
	}
	right := r.Context().Value("rights")
	if right == "auth" {
		tmpl.Execute(w, Auth{IsAuth: true})
		return
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		w.Write([]byte("Плохой запрос"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

type Auth struct {
	IsAuth bool
}

func (h *Handler) Private(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(RenderTemplate("private.html")...)
	if err != nil {
		log.Printf("Parse template private error: %v", err)
	}
	tmpl.ExecuteTemplate(w, "private.html", Auth{IsAuth: true})
}
