// Simplified testharness.js for Sobek runtime
// This is a minimal implementation that provides the essential WPT test functions
// without browser-specific dependencies.

/**
 * Create a synchronous test
 *
 * @param {Function} func - Test function. This is executed
 * immediately. If it returns without error, the test status is
 * set to PASS. If it throws an AssertionError, or
 * any other exception, the test status is set to FAIL.
 * @param {String} name - Test name.
 */
function test(func, name, properties) {
    try {
        func();
    } catch (e) {
        if (name) {
            throw new Error(`Test "${name}" failed: ${e}`);
        }
        throw e;
    }
}

/**
 * Create a promise test (simplified - runs synchronously in Sobek)
 */
function promise_test(func, name, properties) {
    // In Sobek, we run promise tests synchronously
    // This is a simplified implementation
    try {
        const result = func();
        if (result && typeof result.then === 'function') {
            // For actual async tests, we'd need proper async handling
            // For now, this is a placeholder
        }
    } catch (e) {
        if (name) {
            throw new Error(`Promise test "${name}" failed: ${e}`);
        }
        throw e;
    }
}

/**
 * Assert that actual is strictly true
 *
 * @param {Any} actual - Value that is asserted to be true
 * @param {string} [description] - Description of the condition being tested
 */
function assert_true(actual, description) {
    if (!actual) {
        throw new Error(`assert_true: ${description || ''} expected true got ${actual}`);
    }
}

/**
 * Assert that actual is strictly false
 *
 * @param {Any} actual - Value that is asserted to be false
 * @param {string} [description] - Description of the condition being tested
 */
function assert_false(actual, description) {
    if (actual) {
        throw new Error(`assert_false: ${description || ''} expected false got ${actual}`);
    }
}

/**
 * Check if two values are the same (handles NaN and -0)
 */
function same_value(x, y) {
    if (y !== y) {
        // NaN case
        return x !== x;
    }
    if (x === 0 && y === 0) {
        // Distinguish +0 and -0
        return 1/x === 1/y;
    }
    return x === y;
}

/**
 * Assert that actual is the same value as expected.
 *
 * @param {Any} actual - Test value.
 * @param {Any} expected - Expected value.
 * @param {string} [description] - Description of the condition being tested.
 */
function assert_equals(actual, expected, description) {
    if (typeof actual !== typeof expected) {
        throw new Error(`assert_equals: ${description || ''} expected (${typeof expected}) ${expected} but got (${typeof actual}) ${actual}`);
    }
    
    if (!same_value(actual, expected)) {
        throw new Error(`assert_equals: ${description || ''} expected (${typeof expected}) ${expected} but got (${typeof actual}) ${actual}`);
    }
}

/**
 * Assert that actual is not the same value as expected.
 *
 * @param {Any} actual - Test value.
 * @param {Any} expected - The value actual is expected to be different to.
 * @param {string} [description] - Description of the condition being tested.
 */
function assert_not_equals(actual, expected, description) {
    if (same_value(actual, expected)) {
        throw new Error(`assert_not_equals: ${description || ''} got disallowed value ${actual}`);
    }
}

/**
 * Assert that expected is an array and actual is one of the members.
 *
 * @param {Any} actual - Test value.
 * @param {Array} expected - An array that actual is expected to be a member of.
 * @param {string} [description] - Description of the condition being tested.
 */
function assert_in_array(actual, expected, description) {
    if (expected.indexOf(actual) === -1) {
        throw new Error(`assert_in_array: ${description || ''} value ${actual} not in array ${expected}`);
    }
}

/**
 * Assert that actual and expected are both arrays with the same values.
 *
 * @param {Array} actual - Test array.
 * @param {Array} expected - Array that is expected to contain the same values as actual.
 * @param {string} [description] - Description of the condition being tested.
 */
function assert_array_equals(actual, expected, description) {
    if (!Array.isArray(actual)) {
        throw new Error(`assert_array_equals: ${description || ''} actual is not an array`);
    }
    if (!Array.isArray(expected)) {
        throw new Error(`assert_array_equals: ${description || ''} expected is not an array`);
    }
    if (actual.length !== expected.length) {
        throw new Error(`assert_array_equals: ${description || ''} lengths differ, expected ${expected.length} got ${actual.length}`);
    }
    for (let i = 0; i < actual.length; i++) {
        if (!same_value(actual[i], expected[i])) {
            throw new Error(`assert_array_equals: ${description || ''} at index ${i} expected ${expected[i]} got ${actual[i]}`);
        }
    }
}

/**
 * Assert a JS Error with the expected constructor is thrown.
 *
 * @param {Function} constructor - The expected exception constructor.
 * @param {Function} func - Function which should throw.
 * @param {string} [description] - Error description for the case that the error is not thrown.
 */
function assert_throws_js(constructor, func, description) {
    try {
        func();
    } catch (e) {
        if (e instanceof constructor) {
            return;
        }
        // Check by name if instanceof fails (cross-realm issues)
        if (e && e.constructor && e.constructor.name === constructor.name) {
            return;
        }
        throw new Error(`assert_throws_js: ${description || ''} expected ${constructor.name} but got ${e.name || e.constructor?.name || typeof e}: ${e.message || e}`);
    }
    throw new Error(`assert_throws_js: ${description || ''} expected ${constructor.name} but no exception was thrown`);
}

/**
 * Assert a DOMException with the expected type is thrown.
 *
 * @param {number|string} type - The expected exception name or code.
 * @param {Function} func - Function which should throw.
 * @param {string} [description] - Error description.
 */
function assert_throws_dom(type, func, description) {
    try {
        func();
    } catch (e) {
        if (e && (e.name === type || e.code === type)) {
            return;
        }
        throw new Error(`assert_throws_dom: ${description || ''} expected ${type} but got ${e.name || e.code}`);
    }
    throw new Error(`assert_throws_dom: ${description || ''} expected ${type} but no exception was thrown`);
}

/**
 * Asserts if called. Used to ensure that a specific codepath is not taken.
 *
 * @param {string} [description] - Description of the condition being tested.
 */
function assert_unreached(description) {
    throw new Error(`assert_unreached: reached unreachable code${description ? ', reason: ' + description : ''}`);
}

/**
 * Assert that a class string is correct
 */
function assert_class_string(object, expected, description) {
    const actual = Object.prototype.toString.call(object);
    const expectedStr = `[object ${expected}]`;
    if (actual !== expectedStr) {
        throw new Error(`assert_class_string: ${description || ''} expected ${expectedStr} but got ${actual}`);
    }
}

/**
 * Assert that object has own property
 */
function assert_own_property(object, property, description) {
    if (!object.hasOwnProperty(property)) {
        throw new Error(`assert_own_property: ${description || ''} expected property ${property}`);
    }
}

/**
 * Assert that object does not have own property
 */
function assert_not_own_property(object, property, description) {
    if (object.hasOwnProperty(property)) {
        throw new Error(`assert_not_own_property: ${description || ''} unexpected property ${property}`);
    }
}

/**
 * Assert that object inherits from expected
 */
function assert_inherits(object, property, description) {
    if (!(property in object) || object.hasOwnProperty(property)) {
        throw new Error(`assert_inherits: ${description || ''} expected inherited property ${property}`);
    }
}

/**
 * Assert that object has idl attribute
 */
function assert_idl_attribute(object, attribute, description) {
    if (!(attribute in object)) {
        throw new Error(`assert_idl_attribute: ${description || ''} expected attribute ${attribute}`);
    }
}

/**
 * Assert readonly property
 */
function assert_readonly(object, property, description) {
    const desc = Object.getOwnPropertyDescriptor(object, property);
    if (!desc || desc.writable !== false) {
        throw new Error(`assert_readonly: ${description || ''} expected ${property} to be readonly`);
    }
}

/**
 * Assert that a regexp matches
 */
function assert_regexp_match(actual, expected, description) {
    if (!expected.test(actual)) {
        throw new Error(`assert_regexp_match: ${description || ''} expected ${expected} to match ${actual}`);
    }
}

/**
 * Format a value for display in error messages
 */
function format_value(value) {
    if (typeof value === 'string') {
        return JSON.stringify(value);
    }
    if (value === null) {
        return 'null';
    }
    if (value === undefined) {
        return 'undefined';
    }
    if (Array.isArray(value)) {
        return '[' + value.map(format_value).join(', ') + ']';
    }
    if (typeof value === 'object') {
        try {
            return JSON.stringify(value);
        } catch (e) {
            return String(value);
        }
    }
    return String(value);
}

// Setup function (no-op in Sobek)
function setup(func_or_properties, maybe_properties) {
    if (typeof func_or_properties === 'function') {
        func_or_properties();
    }
}

// Done function (no-op in Sobek)
function done() {}

// Add completion callback (no-op in Sobek)
function add_completion_callback(callback) {}

// Add result callback (no-op in Sobek)
function add_result_callback(callback) {}

// Add start callback (no-op in Sobek)
function add_start_callback(callback) {}
