var SecuritySchemes = function(){
    var self = this;

    self.name = ko.observable("");
}

var MethodResponse = function(){
    var self = this;

    self.code = ko.observable(null);
    // self.body = ko.observable(null);
}

var FormRaml = function(){
  var self = this;

  self.schemaSource = ko.observableArray([]);
  self.httpResponseMaster = ko.observableArray([]);

  ajaxPost("/masterhttpstatuses/getallhttpstatuses", { }, function(res){
    self.httpResponseMaster(res.Data.map(function(v){
      return v.code;
    }));
  })

  self.filename = ko.observable("");
  self.title = ko.observable("");
  self.resourceName = ko.observable("");
  self.description = ko.observable("");
  self.selectedSchema = ko.observableArray([]);
  self.securitySchemes = ko.observableArray([new SecuritySchemes()]);
  self.traceability = ko.observable("");
  self.isMethodGet = ko.observable(false);
  self.isMethodPost = ko.observable(false);
  self.methodGetResponses = ko.observableArray([new MethodResponse()]);
  self.methodPostResponses = ko.observableArray([new MethodResponse()]);

  self.getSchemaSource = function(){
    ajaxPost("/masterschema/getallschema", {}, function(res){
      nr.schemaSource(res.Data.map(function(v){
        return {
          text: v.name,
          value: {
            name: v.name,
            body: JSON.parse(v.jsonstring)
          }
        }
      }));
    });
  }

  self.addSecurityScheme = function(){
    self.securitySchemes.push(new SecuritySchemes());
  };

  self.addGetResponse = function(){
    self.methodGetResponses.push(new MethodResponse());
  };

  self.addPostResponse = function(){
    self.methodPostResponses.push(new MethodResponse());
  };

  self.saveRaml = function(){
    self.parent = selectedResource().parent
    self.id = selectedResource()._id

    ajaxPost("/designer/newraml", self, function(res){
      $("#modal-newraml").modal('hide');
    })
  }

  self.reset = function(){
    self.filename("");
    self.title("");
    self.resourceName("");
    self.description("");
    self.selectedSchema([]);
    self.securitySchemes([]);
    self.traceability("");
    self.isMethodGet(false);
    self.isMethodPost(false);
    self.methodGetResponses([new MethodResponse()]);
    self.methodPostResponses([new MethodResponse()]);
  }
}

var FormSchema = function(){
  var self = this;

  self.schname = ko.observable("");
  self.schbody = ko.observable("");

  self.save = function(){
    ajaxPost("/masterschema/save", { id: "", name: self.schname, jsonstring: self.schbody }, function(res){
      nr.getSchemaSource();
      $("#modal-newsch").modal('hide');
    });
  }

  self.reset = function(){
    self.schname("");
    self.schbody("");
  }
}

var nr, sch;

$(document).ready(function(){
  nr = new FormRaml();
  sch = new FormSchema();
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
