ko.validation.init({
  grouping: {
    deep: true,
    live: true,
    observable: true
  }
});

var SecurityResponse = function(){
    var self = this;

    self.code = ko.observable(null);
    self.description = ko.observable("");
}

var SecurityBody = function(){
  var self = this;

  self.responses = ko.observableArray([]);
  self.authorizationUri = ko.observable("");
  self.accessTokenUri = ko.observable("");
  self.authorizationGrants = ko.observable("");

  self.addResponse = function(){
    nr.currentSecuritySchemes().currentBody().responses.push(new SecurityResponse());
  };

  self.removeElement = function(data){
    nr.currentSecuritySchemes().currentBody().responses.remove(data);
  };
}

var SecuritySchemes = function(){
  var self = this;

  self.securitySchemesSource = ko.observableArray(["OAuth 1.0", "OAuth 2.0", "x-{other}"]);

  self.name = ko.observable("").extend({ required: true });
  self.type = ko.observable(null);
  self.description = ko.observable("");
  self.body = ko.observable(null);
  self.currentBody = ko.observable(null);

  self.showDetails = function(){
    nr.currentSecuritySchemes(self);

    if(self.body() == null) {
      self.currentBody(new SecurityBody());
    } else {
      self.currentBody(ko.mapping.fromJS(ko.toJS(self.body)));
    }

    $("#modal-oauth2").modal("show");
  }

  self.changed = function(e){
    if(e.dataItem == "OAuth 2.0"){
      self.showDetails();
    } else {
      self.body(null);
    }
  };

  self.save = function(){
    self.body(self.currentBody());

    $("#modal-oauth2").modal('hide');
  }
}

var MethodTraceability = function(){
    var self = this;

    self.name = ko.observable("");
    self.body = ko.observable("");
}

var MethodResponse = function(){
    var self = this;

    self.code = ko.observable(null);
    self.body = ko.observable(null);
}

var FormRaml = function(){
  var self = this;

  self.schemaSource = ko.observableArray([]);
  self.httpResponseSource = ko.observableArray([]);

  ajaxPost("/masterhttpstatuses/getallhttpstatuses", { }, function(res){
    self.httpResponseSource(res.Data.map(function(v){
      return v.code;
    }));
  });

  self.currentSsInit = function(){
    var ss;

    if(self.securitySchemes().length > 0){
      ss = self.securitySchemes()[0];
    } else {
      ss = new SecuritySchemes();

      Object.keys(ss).forEach(function(name) {
        if (ko.validation.utils.isValidatable(ss[name])) {
          ss[name].rules.removeAll();
        }
      });
    };

    return ss;
  }

  self.filename = ko.observable("").extend({ required: true });
  self.resourceName = ko.observable("");
  self.description = ko.observable("");
  self.selectedSchema = ko.observableArray([]);
  self.securitySchemes = ko.observableArray([]);
  self.currentSecuritySchemes = ko.observable(self.currentSsInit());
  self.isMethodGet = ko.observable(false);
  self.isMethodPost = ko.observable(false);
  self.methodGetTraceabilities = ko.observableArray([]);
  self.methodPostTraceabilities = ko.observableArray([]);
  self.methodGetResponses = ko.observableArray([]);
  self.methodPostResponses = ko.observableArray([]);

  self.traceabilitySource = ko.computed(function(){
    var standardType = ["string", "boolean", "number", "integer", "date-only", "time-only", "datetime-only", "datetime", "file"];
    return self.selectedSchema().map(function(v){
      return v.name;
    }).concat(standardType);
  });

  self.getSchemaSource = function(){
    ajaxPost("/masterschema/getallschema", {}, function(res){
      nr.schemaSource(res.Data.map(function(v){
        return {
          text: v.name,
          value: {
            name: v.name,
            body: jsonAttrStringify(v.jsonstring)
          }
        }
      }));
    });
  }

  self.addSecurityScheme = function(){
    self.securitySchemes.push(new SecuritySchemes());
  };

  self.addGetTraceability = function(){
    self.methodGetTraceabilities.push(new MethodTraceability());
  };

  self.addPostTraceability = function(){
    self.methodPostTraceabilities.push(new MethodTraceability());
  };

  self.addGetResponse = function(){
    self.methodGetResponses.push(new MethodResponse());
  };

  self.addPostResponse = function(){
    self.methodPostResponses.push(new MethodResponse());
  };

  self.removeElement = function(obs){
    return function(data){
      self[obs].remove(data);
    }
  }

  self.validate = function(){
    if (self.errors().length === 0) {
      return true;
    } else {
      self.errors.showAllMessages();
      return false;
    }
  }

  self.saveRaml = function(){
    self.parent = selectedResource().parent;
    self.id = selectedResource()._id;
    self.baseUri = selectedResource().uri;

    if(self.validate()){
      ajaxPost("/designer/newraml", self, function(res){
        $("#modal-newraml").modal('hide');

        if(stsTabForm() == 'raml'){
          refreshEditor();
        }
      })
    }
  }

  self.reset = function(){
    self.filename("");
    self.resourceName("");
    self.description("");
    self.selectedSchema([]);
    self.securitySchemes([]);
    self.currentSecuritySchemes(self.currentSsInit());
    self.isMethodGet(false);
    self.isMethodPost(false);

    self.methodGetTraceabilities([]);
    self.methodPostTraceabilities([]);
    self.methodGetResponses([]);
    self.methodPostResponses([]);

    self.errors.showAllMessages(false);
  }
}

var FormSchema = function(){
  var self = this;

  self.schname = ko.observable("").extend({ required: true });
  self.schbody = ko.observable("")
                  .extend({ required: true })
                  .extend({
                    validation: {
                      validator: IsJsonString,
                      message: 'JSON is not valid.'
                    }
                  });

  self.validate = function(){
    if (self.errors().length === 0) {
      return true;
    } else {
      self.errors.showAllMessages();
      return false;
    }
  }

  self.save = function(){
    if(self.validate()){
      ajaxPost("/masterschema/save", { id: "", name: self.schname, jsonstring: self.schbody }, function(res){
        nr.getSchemaSource();
        $("#modal-newsch").modal('hide');
      });
    }
  }

  self.reset = function(){
    self.schname("");
    self.schbody("");

    self.errors.showAllMessages(false);
  }
}

var nr, sch;

$(document).ready(function(){
  nr = new FormRaml();
  sch = new FormSchema();

  nr.errors = ko.validation.group(nr);
  sch.errors = ko.validation.group(sch);
})

$('#modal-newraml').on('hidden.bs.modal', function () {
  nr.reset();
});

$('#modal-newraml').on('shown.bs.modal', function (e) {
  ko.bindingHandlers.kendoMultiSelect.options.filter = "contains";

  nr.getSchemaSource();
});

$('#modal-newsch').on('hidden.bs.modal', function () {
  sch.reset();
});

$(document).on('show.bs.modal', '.modal', function (event) {
  var zIndex = 1040 + (10 * $('.modal:visible').length);
  $(this).css('z-index', zIndex);
  setTimeout(function() {
    $('.modal-backdrop').not('.modal-stack').css('z-index', zIndex - 1).addClass('modal-stack');
  }, 0);
});
