var ansiconf = {
    standalone: true,
    escapeHtml: true,
    palette: {
    black: '#0f3059',
    white: '#eeeeee',
    red:   '#dd0000',
    green: '#00cc3e',
    blue:  '#0099ff',
    yellow:'#eeee00',
    purple:'#bb00bb',
    cyan:  '#eeeeee'
}
}

/**
 * Helper for escaping HTML entities; extracted from lodash codebase
 * @param {String} string
 */
function escapeHtml(string) {
    var reUnescapedHtml = /[&<>"'`]/g;
    var reHasUnescapedHtml = RegExp(reUnescapedHtml.source);
    var htmlEscapes = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;',
        '`': '&#96;'
    };
    var escapeHtmlChar = function (chr) {
        return htmlEscapes[chr];
    };
    return (string && reHasUnescapedHtml.test(string)) ?
        string.replace(reUnescapedHtml, escapeHtmlChar) : string;
}

/**
 * Generate HTML header including CSS style with palette
 * @param {Object} palette The following keys are accepted:
 *   fg, fg_{black,red,green,yellow,blue,purple,cyan,white},
 *   bg, bg_{black,red,green,yellow,blue,purple,cyan,white},
 *           black,red,green,yellow,blue,purple,cyan,white
 * Each of them is optional and defaults to its named HTML correspondent
 */
var makeHeader = function (palette) {
    var p = palette || {};
    var str = '<!DOCTYPE html><html><head><meta charset="utf-8"><style>\
    .ansi_bold{font-weight:bold}\
    .ansi_italic{font-style:italic}\
    .ansi_console_snippet{font-family:monospace; white-space: pre; display: block;\
      unicode-bidi: embed; overflow:initial;word-wrap: break-word; padding:5px;}';

    str += '\
    .ansi_console_snippet{\
        color:'           +(p.fg || p.fg_white || p.white || "white")+';\
    }\
    .ansi_fg_black  {color:'+(p.fg_black  || p.black || "black ")+'}\
    .ansi_fg_red    {color:'+(p.fg_red    || p.red   || "red   ")+'}\
    .ansi_fg_green  {color:'+(p.fg_green  || p.green || "green ")+'}\
    .ansi_fg_yellow {color:'+(p.fg_yellow || p.yellow|| "yellow")+'}\
    .ansi_fg_blue   {color:'+(p.fg_blue   || p.blue  || "blue  ")+'}\
    .ansi_fg_purple {color:'+(p.fg_purple || p.purple|| "purple")+'}\
    .ansi_fg_cyan   {color:'+(p.fg_cyan   || p.cyan  || "cyan  ")+'}\
    .ansi_fg_white  {color:'+(p.fg_white  || p.white || "white ")+'}\
    ';

    str += '</style>';
    return str;
};

var headerLocal = '<span class="ansi_console_snippet">';
var footerLocal = '</span>';
var footer = '</body></html>';

/**
 * Internal method which escapes HTML entities if requested, and converts shell escape codes
 * to HTML `<span class=...>` markers
 * @param {String} inputString
 * @param {Object} options
 * @return {String}
 */
function processString (inputString, options) {
    inputString = options.escapeHtml ? escapeHtml(inputString) : inputString;
    var str = ansispan(inputString);
    while (true) {
        var openSpanCount = (str.match(/<span/g) || []).length;
        var closeSpanCount = (str.match(/<\/span/g) || []).length;
        if (openSpanCount >= closeSpanCount) {
            break;
        }
        // we may have some extra closing escape sequence that doesn't really close anything
        // especially on end of line but not only
        var underReplace = '<\/span>';
        var idx = str.lastIndexOf(underReplace);
        if (idx > -1) {
            str = str.substring(0, idx) + str.substring(idx + underReplace.length);
        }
    }
    return str;
}

/**
 * Helper to normalize config options
 * @param {Object} options, defaults to { standalone: true, escapeHtml: true, palette: defaultPalette }
 * @return {Object}
 */
function processOptions (options) {
    options = options || {};
    if (typeof options.escapeHtml === "undefined") {
        options.escapeHtml = true;
    }
    if (typeof options.wrapped === "undefined") {
        options.wrapped = options.standalone;
    }

    if (!options.standalone && options.palette) {
        console.error('[ansi2html] options.standalone == false; palette will be ignored');
    }
    return options;
}

/**
 * Convert the text from `inputStream` and output it to `outputStream`
 * @param {Object} options
 * @param {Stream} inputStream defaults to STDIN
 * @param {Stream} outputStream defaults to STDOUT
 */
function ansi2html_stream (options, inputStream, outputStream) {
    inputStream = inputStream || process.stdin;
    outputStream = outputStream || process.stdout;

    function write(str) { outputStream.write(str); }
    var input = readline.createInterface({
        input: inputStream,
        output: outputStream
    });

    if (typeof options.standalone === "undefined") {
        options.standalone = true;
    }
    /////////////////////////////////////////////

    options = processOptions(options);
    if (options.standalone) write(makeHeader(options.palette));
    if (options.wrapped)    write(headerLocal);
    input.on('line', function(line){
        write(processString(line, options) + EOL);
    });
    input.on('close', function() {
        if (options.wrapped)    write(footerLocal);
        if (options.standalone) write(footer);
    });
}

/**
 * Convert the text from `inputString` and return it
 * @param {Object} options
 * @return {String}
 */
function ansi2html_string (options, inputString) {
    if (arguments.length < 2) {
        // support a2h(inputString) format
        inputString = options;
        options = {};
    }

    var out = "";
    function write(str) { out += str; }

    if (typeof options.standalone === "undefined") {
        options.standalone = false;
    }
    /////////////////////////////////////////////

    options = processOptions(options);
    if (options.standalone) write(makeHeader(options.palette));
    if (options.wrapped)    write(headerLocal);
    write(processString(inputString, options));
    if (options.wrapped)    write(footerLocal);
    if (options.standalone) write(footer);

    return out;
}

var ansispan = function (str) {

    //
    // `\033[Xm` == `\033[0;Xm` sets foreground color to `X`.
    //
    str = str.replace(
        /(\033\[(\d+)(;\d+)?m)/gm,
        function(match, fullMatch, m1, m2) {
            var fgColor = m1;
            var bgColor = m2;

            var newStr = '<span class="';
            if (fgColor && ansispan.foregroundColors[fgColor]) {
                newStr += 'ansi_fg_' + ansispan.foregroundColors[fgColor];
            }
            if (bgColor) {
                bgColor = bgColor.substr(1); // remove leading ;
                if (ansispan.backgroundColors[bgColor]) {
                    newStr += ' ansi_bg_' + ansispan.backgroundColors[bgColor];
                }
            }
            newStr += '">';
            return newStr;
        }
    );

    //
    // `\033[1m` enables bold font, `\033[22m` disables it
    //
    str = str.replace(/\033\[1m/g, '<span class="ansi_bold">').replace(/\033\[22m/g, '</span>');

    //
    // `\033[3m` enables italics font, `\033[23m` disables it
    //
    str = str.replace(/\033\[3m/g, '<span class="ansi_italic">').replace(/\033\[23m/g, '</span>');

    str = str.replace(/\033\[m/g, '</span>');
    str = str.replace(/\033\[0m/g, '</span>');
    return str.replace(/\033\[39m/g, '</span>');
};

ansispan.foregroundColors = {
    '30': 'black',
    '31': 'red',
    '32': 'green',
    '33': 'yellow',
    '34': 'blue',
    '35': 'purple',
    '36': 'cyan',
    '37': 'white'
};

ansispan.backgroundColors = {
    '40': 'black',
    '41': 'red',
    '42': 'green',
    '43': 'yellow',
    '44': 'blue',
    '45': 'purple',
    '46': 'cyan',
    '47': 'white'
};