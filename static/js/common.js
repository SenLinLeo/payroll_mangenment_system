/**
 * 功能：将对象转化为Key-Value键值对形式的字符串
 * 传入参数：对象名
 * 返回参数：Key-Value字符串
 **/
function objectToKvpString(object, srcDict = new Array()) {
    var flag = 0
        , strKey = "";
    for (var key in object) {
        if (object[key] == "" && srcDict[key] != "") {
            object[key] = srcDict[key];
        }
        if (0 == flag) {
            flag = 1;
            strKey = key + '=' + object[key];

        } else {
            strKey = strKey + "&" + key + '=' + object[key];
        }
    }

    return strKey;
}