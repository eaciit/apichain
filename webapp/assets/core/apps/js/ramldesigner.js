var filesUrl = "/files"

var ramldesigner = {}
ramldesigner.ramlAddr = ko.observable(filesUrl);
ramldesigner.buttondetailvis = ko.observable(false);
ramldesigner.filename = function(){
  return $("#ramlapp > raml-editor > div > ul > li.file-absolute-path").text().split(" ")[0];
};

ramldesigner.resource = ko.observable("");
ramldesigner._resource = ko.observable("");

ramldesigner.tabs = ko.observableArray([]);

function myFileSystem($http, $q, config, $rootScope) {
  var service = {};
  var handler = "/designer"

  service.directory = function (path) {
    var deferred = $q.defer();

    ajaxPost(handler + "/list", { parent: selectedResource().parent(), id: selectedResource()._id() }, function(res){
      if(res.success){
        $rootScope.$broadcast('event:notification', {
          message : 'Directory loaded.',
          expires : true
        });
        deferred.resolve(res.data);
      } else {
        var err = (res.message || status || 'Unknown Error');
        $rootScope.$broadcast('event:notification', {
          message : 'Directory NOT loaded: ' + err,
          level : 'error',
          expires : false
        });
        deferred.reject.bind(deferred);
      }
    });

    return deferred.promise;
  };

  service.load = function (path, name) {
    ramldesigner.buttondetailvis(false);

    var deferred = $q.defer();
    if (path.endsWith(".meta")) {
      deferred.resolve("{}");
      return deferred.promise;
    }

    var succeeded = function(data, status, headers, config){
      $rootScope.$broadcast('event:notification', {
        message : 'File loaded.',
        expires : true
      });

      ramldesigner.buttondetailvis(true);
      deferred.resolve(data);
    };

    var failed = function(data, status, headers, config){
      var err = (data.error || status || 'Unknown Error');
      $rootScope.$broadcast('event:notification', {
        message : 'File NOT loaded: ' + err,
        level : 'error',
        expires : false
      });
      deferred.reject.bind(deferred);
    };

    ajaxGET(ramldesigner.ramlAddr() + path, succeeded, failed);

    return deferred.promise;
  };

  service.remove = function (path, name) {
    var deferred = $q.defer();

    ajaxPost(handler + "/delete", { parent: selectedResource().parent(), id: selectedResource()._id(), path: path, name: name }, function(res){
      if(res.success){
        $rootScope.$broadcast('event:notification', {
          message : 'File removed.',
          expires : true
        });
        deferred.resolve();
      } else {
        var err = (res.message || status || 'Unknown Error');
        $rootScope.$broadcast('event:notification', {
          message : 'File NOT removed: ' + err,
          level : 'error',
          expires : false
        });
        deferred.reject.bind(deferred);
      }
    });

    return deferred.promise;
  };

  service.save = function (path, contents) {
    var deferred = $q.defer();
    if (path.endsWith(".meta")) {
      deferred.resolve();
      return deferred.promise;
    }

    ajaxPost(handler + "/save", { parent: selectedResource().parent(), id: selectedResource()._id(), path: path, contents: contents }, function(res){
      if(res.success){
        $rootScope.$broadcast('event:notification', {
          message : 'File saved.',
          expires : true
        });
        deferred.resolve();
      } else {
        var err = (res.message || status || 'Unknown Error');
        $rootScope.$broadcast('event:notification', {
          message : 'File NOT saved: ' + err,
          level : 'error',
          expires : false
        });
        deferred.reject.bind(deferred);
      }
    });

    return deferred.promise;
  };

  service.createFolder = function(path) {
		var deferred = $q.defer();

    ajaxPost(handler + "/newfolder", { parent: selectedResource().parent(), id: selectedResource()._id(), path: path }, function(res){
      if(res.success){
  			$rootScope.$broadcast('event:notification', {
  				message : 'Folder created.',
  				expires : true
  			});
        deferred.resolve();
      } else {
        var err = (res.message || status || 'Unknown Error');
  			$rootScope.$broadcast('event:notification', {
  				message : 'Folder NOT created: ' + err,
  				level : 'error',
  				expires : false
  			});
        deferred.reject.bind(deferred);
      }
    });

		return deferred.promise;
	};

  service.rename = function(source, destination) {
		var deferred = $q.defer();

    ajaxPost(handler + "/rename", { parent: selectedResource().parent(), id: selectedResource()._id(), source: source, destination: destination }, function(res){
      if(res.success){
  			$rootScope.$broadcast('event:notification', {
  				message : 'File renamed.',
  				expires : true
  			});
  			deferred.resolve();
      } else {
        var err = (res.message || status || 'Unknown Error');
  			$rootScope.$broadcast('event:notification', {
  				message : 'File NOT renamed: ' + err,
  				level : 'error',
  				expires : false
  			});
  			deferred.reject.bind(deferred);
      }
    });

		return deferred.promise;
	};

  return service;
}

var stripText = function(text){
  return text.replace(/\s/g, '');
}

var getLines = function(){
  var $lines = $(".CodeMirror-code > div")
  var ret = []

  $lines.each(function(i, v){
  	var $line = $(v);
  	var level = $line.find(".cm-indent").length;
  	var method = $line.find(".cm-method-content").text();
  	var content = $line.find("pre > span").text();

  	var $lineBefore = (i == 0) ? undefined : $($lines[i - 1]);
  	var levelBefore = $lineBefore ? $lineBefore.find(".cm-indent").length : undefined;

  	var line = {
  		lineNum: i + 1,
  		level: level,
  		method: method,
  		content: content,
  	};

  	ret.push(line);
  })

  return ret;
}

var bindRaml = function(){
  stsTabForm("raml");
  syncform();

  setTimeout(function(){
    refreshEditor();

    $('#modal-ramldet').on('shown.bs.modal', function (e) {
    });

    $('#modal-ramldet').on('hidden.bs.modal', function () {
    });
  }, 200);
};

var saveFromForm = function(){
  var lines = getLines();

  console.log(lines);
  var found = lines.find(function(v){
  	return v.method == "description" && v.content.indexOf(ramldesigner.tabs()[0].sections.descriptionBefore().trim()) !== -1
  })

  ajaxPost("/designer/editline", {
    path: ramldesigner.filename(),
    line: found.lineNum,
    to: found.content.replace(ramldesigner.tabs()[0].sections.descriptionBefore().trim(), ramldesigner.tabs()[0].sections.description())
  }, function(res){
  	console.log(res)

    refreshEditor();
  })
}

var refreshEditor = function(){
  $("#ramlapp").html("");
  angular.module('ramlEditorApp')
    .config(function (fileSystemProvider) {
      fileSystemProvider.setFileSystemFactory(myFileSystem);
    });

  $("#ramlapp").replaceWith('<div id="ramlapp" class="container" ng-app="ramlEditorApp"><raml-editor></raml-editor></div>');
  angular.bootstrap($("#ramlapp"), ['ramlEditorApp']);

  $("raml-editor > div > header").remove();
  $(".menu-item-mocking-service").remove();
  $("#consoleAndEditor").remove();

  $(".raml-console-embedded").appendTo("#modal-ramldet .modal-body");
}

$(document).ready(function(){
  $(document).bind('keydown', function(e) {
    if ((e.metaKey || e.ctrlKey) && (String.fromCharCode(e.which).toLowerCase() === 's')) {
      e.preventDefault();
      return false;
    }
  });

  var waitForElement = function(selector, callback) {
    if(document.querySelector(selector)!=null) {
      callback();
      return;
    }
    else {
      setTimeout(function() {
        waitForElement(selector, callback);
      }, 0);
    }
  }

  waitForElement(".modal.fade.ng-isolate-scope.in", function(){
    $(".modal.fade.ng-isolate-scope.in").remove();
    $(".modal-backdrop.fade.in").remove();
  });
})
