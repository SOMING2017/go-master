Vue.http.interceptors.push(function (request) {
    request.headers.set('Content-Type', 'application/json');
});
