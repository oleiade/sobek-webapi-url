// Stubs for browser APIs not available in Sobek runtime
// These are required for running WPT tests

// Stub for DOMException (not available in Sobek)
var DOMException = (function() {
    function DOMException(message, name) {
        this.message = message;
        this.name = name || 'Error';
    }
    DOMException.prototype = Object.create(Error.prototype);
    DOMException.prototype.constructor = DOMException;
    // Add standard error codes
    DOMException.INDEX_SIZE_ERR = 1;
    DOMException.DOMSTRING_SIZE_ERR = 2;
    DOMException.HIERARCHY_REQUEST_ERR = 3;
    DOMException.WRONG_DOCUMENT_ERR = 4;
    DOMException.INVALID_CHARACTER_ERR = 5;
    DOMException.NO_DATA_ALLOWED_ERR = 6;
    DOMException.NO_MODIFICATION_ALLOWED_ERR = 7;
    DOMException.NOT_FOUND_ERR = 8;
    DOMException.NOT_SUPPORTED_ERR = 9;
    DOMException.INUSE_ATTRIBUTE_ERR = 10;
    DOMException.INVALID_STATE_ERR = 11;
    DOMException.SYNTAX_ERR = 12;
    DOMException.INVALID_MODIFICATION_ERR = 13;
    DOMException.NAMESPACE_ERR = 14;
    DOMException.INVALID_ACCESS_ERR = 15;
    DOMException.VALIDATION_ERR = 16;
    DOMException.TYPE_MISMATCH_ERR = 17;
    DOMException.SECURITY_ERR = 18;
    DOMException.NETWORK_ERR = 19;
    DOMException.ABORT_ERR = 20;
    DOMException.URL_MISMATCH_ERR = 21;
    DOMException.QUOTA_EXCEEDED_ERR = 22;
    DOMException.TIMEOUT_ERR = 23;
    DOMException.INVALID_NODE_TYPE_ERR = 24;
    DOMException.DATA_CLONE_ERR = 25;
    return DOMException;
})();

// Stub for FormData (not available in Sobek)
var FormData = (function() {
    function FormData() {
        this._entries = [];
    }
    FormData.prototype.append = function(name, value) {
        this._entries.push([name, value]);
    };
    FormData.prototype.get = function(name) {
        for (var i = 0; i < this._entries.length; i++) {
            if (this._entries[i][0] === name) {
                return this._entries[i][1];
            }
        }
        return null;
    };
    FormData.prototype.has = function(name) {
        for (var i = 0; i < this._entries.length; i++) {
            if (this._entries[i][0] === name) {
                return true;
            }
        }
        return false;
    };
    FormData.prototype[Symbol.iterator] = function() {
        return this._entries[Symbol.iterator]();
    };
    return FormData;
})();

/**
 * Run tests by key subset (simplified - just runs all tests)
 */
function subsetTestByKey(key, testFunc, func, name) {
    testFunc(func, name);
}

