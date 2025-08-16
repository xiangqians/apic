window.requestInterceptor = (function () {
    let interceptor = {};

    interceptor.handle = function (operation, request) {
        console.log(operation, request);
        return request;
    }

    return interceptor;
})();