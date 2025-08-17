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

/**
 * 扩展日期格式化函数
 * yyyy/MM/dd HH:mm:ss
 * yyyy/MM/dd HH:mm:ss.SSS
 *
 * @param pattern
 * @returns {string}
 */
Date.prototype.format = function (pattern) {
    let object = {
        "yyyy": this.getFullYear().toString(),
        "yy": this.getFullYear().toString().substring(2),
        "MM": (this.getMonth() + 1).toString().padStart(2, '0'),
        "dd": this.getDate().toString().padStart(2, '0'),
        "HH": this.getHours().toString().padStart(2, '0'),
        "mm": this.getMinutes().toString().padStart(2, '0'),
        "ss": this.getSeconds().toString().padStart(2, '0'),
        "SSS": this.getMilliseconds().toString().padStart(3, '0'),
    };
    return pattern.replace(/yyyy|yy|MM|dd|HH|mm|ss|SSS/g, match => object[match]);
}