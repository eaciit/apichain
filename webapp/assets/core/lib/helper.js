var multipleAjaxPost = function(ajaxposts, callback) {
    var doRequest = function(){
        var deferred = $.Deferred();

        $.when(...arguments).done(function () {
            if (ajaxposts.length == 1)
                deferred.resolve([arguments])
            else
                deferred.resolve(arguments)
        });

        return deferred.promise();
    }

    doRequest(...ajaxposts).then(callback);
}

var findValuesHelper = function(obj, key, list) {
  if (!obj) return list;
  if (obj instanceof Array) {
    for (var i in obj) {
        list = list.concat(findValuesHelper(obj[i], key, []));
    }
    return list;
  }
  if (obj[key]) list.push(obj[key]);

  if ((typeof obj == "object") && (obj !== null) ){
	  var children = Object.keys(obj);
	  if (children.length > 0){
	  	for (i = 0; i < children.length; i++ ){
	        list = list.concat(findValuesHelper(obj[children[i]], key, []));
	  	}
	  }
  }
  return list;
}

var findValues = function(obj, key){
	return findValuesHelper(obj, key, []);
}

function cleanArray(actual) {
  var newArray = new Array();
  for (var i = 0; i < actual.length; i++) {
    if (actual[i]) {
      newArray.push(actual[i]);
    }
  }
  return newArray;
}
