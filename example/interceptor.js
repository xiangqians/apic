// 请求拦截器
window.requestInterceptor = (function () {
    let interceptor = {};
    interceptor.handle = function (operation, request) {
        console.log(`[${new Date().format('yyyy/MM/dd HH:mm:ss.SSS')} - request handle]`, operation, request);
        return request;
    }
    return interceptor;
})();

// 响应拦截器
window.responseInterceptor = (function () {
    let interceptor = {};
    interceptor.handle = function (operation, response) {
        console.log(`[${new Date().format('yyyy/MM/dd HH:mm:ss.SSS')} - response handle]`, operation, response);
        return response;
    }
    return interceptor;
})();