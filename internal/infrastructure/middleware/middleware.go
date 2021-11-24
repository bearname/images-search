package middleware

//
//func DecodeCreateUrlRequest(next http.HandlerFunc) http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		log.Println("DecodeCreateUrlRequest")
//		var urlRequest transport.CreateUrlRequest
//		err := json.NewDecoder(request.Body).Decode(&urlRequest)
//		if err != nil {
//			http.Error(writer, "customAlias or originalUrl not present on json", http.StatusBadRequest)
//			return
//		}
//
//		b := !util.IsValidUrl(urlRequest.GetOriginalUrl())
//		if b {
//			http.Error(writer, "invalid original url", http.StatusBadRequest)
//			return
//		}
//
//		context.Set(request, "url", urlRequest)
//
//		next(writer, request)
//	}
//}
