window.requestInterceptor = (function () {
    let requestInterceptor = {};

    requestInterceptor.handle = function (operation, request) {
        console.log(operation, request);
        return request;
    }

    return requestInterceptor;
})();