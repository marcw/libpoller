package check

//import (
	//"net/http"
	//"net/http/httptest"
	//"testing"
//)

//type successPollHandler struct {
//}

//func (p successPollHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
//}

//func TestPollIsContactingHttpServer(t *testing.T) {
	//server := httptest.NewServer(successPollHandler{})
	//defer server.Close()

	//check, _ := NewCheck(server.URL, "foobar", "10s", make(map[string]string))
	//statusCode, _, err := check.Poll()
	//if err != nil {
		//t.Error(err)
	//}
	//if statusCode != 200 {
		//t.Error("statusCode should be 200")
	//}
//}
