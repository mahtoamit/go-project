package constants

const (
	Url_get_books = "/api/v1/books"
	Url_get_single_book = "/api/v1/book/:title"
	Url_add_book = "/api/v1/book/"
	Url_update_book ="/api/v1/book/:id"
	Unmarshal_error = "Error unmarshaling cached books data:"
	Error_caching_data = "Error caching books data:"
	Error_deleting_cached_data = "Error deleting employees cache: "
	Url_login = "/api/v1/login"
	Url_signup = "/api/v1/signup"
	Url_logout = "/api/v1/logout"
	Url_book = "/api/v1/book/:title"
	Redis_book_const = "books:%s"
    Filename = "record_log/app.log"
	Filename_start = "app-"
	Filename_ext = ".log"

)