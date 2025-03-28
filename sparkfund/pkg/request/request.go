package request

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/pkg/errors"
	"github.com/sparkfund/pkg/validator"
)

// BindJSON binds JSON request body to a struct
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := validator.ValidateStruct(obj); err != nil {
		return err
	}
	return nil
}

// BindQuery binds query parameters to a struct
func BindQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := validator.ValidateStruct(obj); err != nil {
		return err
	}
	return nil
}

// BindForm binds form data to a struct
func BindForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := validator.ValidateStruct(obj); err != nil {
		return err
	}
	return nil
}

// BindURI binds URI parameters to a struct
func BindURI(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := validator.ValidateStruct(obj); err != nil {
		return err
	}
	return nil
}

// GetQueryString gets a query parameter as string
func GetQueryString(c *gin.Context, key string) (string, error) {
	value := c.Query(key)
	if value == "" {
		return "", errors.ErrBadRequest(fmt.Errorf("missing required query parameter: %s", key))
	}
	return value, nil
}

// GetQueryInt gets a query parameter as integer
func GetQueryInt(c *gin.Context, key string) (int, error) {
	value := c.Query(key)
	if value == "" {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required query parameter: %s", key))
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid integer value for query parameter: %s", key))
	}
	return intValue, nil
}

// GetQueryFloat gets a query parameter as float
func GetQueryFloat(c *gin.Context, key string) (float64, error) {
	value := c.Query(key)
	if value == "" {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required query parameter: %s", key))
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid float value for query parameter: %s", key))
	}
	return floatValue, nil
}

// GetQueryBool gets a query parameter as boolean
func GetQueryBool(c *gin.Context, key string) (bool, error) {
	value := c.Query(key)
	if value == "" {
		return false, errors.ErrBadRequest(fmt.Errorf("missing required query parameter: %s", key))
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.ErrBadRequest(fmt.Errorf("invalid boolean value for query parameter: %s", key))
	}
	return boolValue, nil
}

// GetFormString gets a form parameter as string
func GetFormString(c *gin.Context, key string) (string, error) {
	value := c.PostForm(key)
	if value == "" {
		return "", errors.ErrBadRequest(fmt.Errorf("missing required form parameter: %s", key))
	}
	return value, nil
}

// GetFormInt gets a form parameter as integer
func GetFormInt(c *gin.Context, key string) (int, error) {
	value := c.PostForm(key)
	if value == "" {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required form parameter: %s", key))
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid integer value for form parameter: %s", key))
	}
	return intValue, nil
}

// GetFormFloat gets a form parameter as float
func GetFormFloat(c *gin.Context, key string) (float64, error) {
	value := c.PostForm(key)
	if value == "" {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required form parameter: %s", key))
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid float value for form parameter: %s", key))
	}
	return floatValue, nil
}

// GetFormBool gets a form parameter as boolean
func GetFormBool(c *gin.Context, key string) (bool, error) {
	value := c.PostForm(key)
	if value == "" {
		return false, errors.ErrBadRequest(fmt.Errorf("missing required form parameter: %s", key))
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.ErrBadRequest(fmt.Errorf("invalid boolean value for form parameter: %s", key))
	}
	return boolValue, nil
}

// GetParamString gets a URI parameter as string
func GetParamString(c *gin.Context, key string) (string, error) {
	value := c.Param(key)
	if value == "" {
		return "", errors.ErrBadRequest(fmt.Errorf("missing required URI parameter: %s", key))
	}
	return value, nil
}

// GetParamInt gets a URI parameter as integer
func GetParamInt(c *gin.Context, key string) (int, error) {
	value := c.Param(key)
	if value == "" {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required URI parameter: %s", key))
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid integer value for URI parameter: %s", key))
	}
	return intValue, nil
}

// GetHeaderString gets a header value as string
func GetHeaderString(c *gin.Context, key string) (string, error) {
	value := c.GetHeader(key)
	if value == "" {
		return "", errors.ErrBadRequest(fmt.Errorf("missing required header: %s", key))
	}
	return value, nil
}

// GetHeaderInt gets a header value as integer
func GetHeaderInt(c *gin.Context, key string) (int, error) {
	value := c.GetHeader(key)
	if value == "" {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required header: %s", key))
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid integer value for header: %s", key))
	}
	return intValue, nil
}

// GetHeaderFloat gets a header value as float
func GetHeaderFloat(c *gin.Context, key string) (float64, error) {
	value := c.GetHeader(key)
	if value == "" {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required header: %s", key))
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid float value for header: %s", key))
	}
	return floatValue, nil
}

// GetHeaderBool gets a header value as boolean
func GetHeaderBool(c *gin.Context, key string) (bool, error) {
	value := c.GetHeader(key)
	if value == "" {
		return false, errors.ErrBadRequest(fmt.Errorf("missing required header: %s", key))
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.ErrBadRequest(fmt.Errorf("invalid boolean value for header: %s", key))
	}
	return boolValue, nil
}

// GetCookieString gets a cookie value as string
func GetCookieString(c *gin.Context, key string) (string, error) {
	value, err := c.Cookie(key)
	if err != nil {
		return "", errors.ErrBadRequest(fmt.Errorf("missing required cookie: %s", key))
	}
	return value, nil
}

// GetCookieInt gets a cookie value as integer
func GetCookieInt(c *gin.Context, key string) (int, error) {
	value, err := c.Cookie(key)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required cookie: %s", key))
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid integer value for cookie: %s", key))
	}
	return intValue, nil
}

// GetCookieFloat gets a cookie value as float
func GetCookieFloat(c *gin.Context, key string) (float64, error) {
	value, err := c.Cookie(key)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("missing required cookie: %s", key))
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.ErrBadRequest(fmt.Errorf("invalid float value for cookie: %s", key))
	}
	return floatValue, nil
}

// GetCookieBool gets a cookie value as boolean
func GetCookieBool(c *gin.Context, key string) (bool, error) {
	value, err := c.Cookie(key)
	if err != nil {
		return false, errors.ErrBadRequest(fmt.Errorf("missing required cookie: %s", key))
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.ErrBadRequest(fmt.Errorf("invalid boolean value for cookie: %s", key))
	}
	return boolValue, nil
}

// GetFile gets a file from the request
func GetFile(c *gin.Context, key string) (*multipart.FileHeader, error) {
	file, err := c.FormFile(key)
	if err != nil {
		return nil, errors.ErrBadRequest(fmt.Errorf("missing required file: %s", key))
	}
	return file, nil
}

// GetFiles gets multiple files from the request
func GetFiles(c *gin.Context, key string) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, errors.ErrBadRequest(fmt.Errorf("invalid multipart form"))
	}
	files := form.File[key]
	if len(files) == 0 {
		return nil, errors.ErrBadRequest(fmt.Errorf("missing required files: %s", key))
	}
	return files, nil
}

// GetJSON gets JSON data from the request body
func GetJSON(c *gin.Context, obj interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := json.Unmarshal(body, obj); err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := validator.ValidateStruct(obj); err != nil {
		return err
	}
	return nil
}

// GetXML gets XML data from the request body
func GetXML(c *gin.Context, obj interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := xml.Unmarshal(body, obj); err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := validator.ValidateStruct(obj); err != nil {
		return err
	}
	return nil
}

// GetYAML gets YAML data from the request body
func GetYAML(c *gin.Context, obj interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := yaml.Unmarshal(body, obj); err != nil {
		return errors.ErrBadRequest(err)
	}
	if err := validator.ValidateStruct(obj); err != nil {
		return err
	}
	return nil
}

// GetRawBody gets raw request body
func GetRawBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.ErrBadRequest(err)
	}
	return body, nil
}

// GetContentType gets the content type of the request
func GetContentType(c *gin.Context) string {
	return c.GetHeader("Content-Type")
}

// IsJSON checks if the request content type is JSON
func IsJSON(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "application/json")
}

// IsXML checks if the request content type is XML
func IsXML(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "application/xml")
}

// IsYAML checks if the request content type is YAML
func IsYAML(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "application/x-yaml")
}

// IsForm checks if the request content type is form
func IsForm(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "application/x-www-form-urlencoded")
}

// IsMultipart checks if the request content type is multipart
func IsMultipart(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "multipart/form-data")
}

// IsText checks if the request content type is text
func IsText(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "text/plain")
}

// IsHTML checks if the request content type is HTML
func IsHTML(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "text/html")
}

// IsBinary checks if the request content type is binary
func IsBinary(c *gin.Context) bool {
	return strings.Contains(GetContentType(c), "application/octet-stream")
} 