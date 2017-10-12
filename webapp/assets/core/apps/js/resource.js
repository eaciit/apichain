var pageProcessing = ko.observable(true)

var stsForm = ko.observable('grid');
var stsExistUri = ko.observable(false);
var stsDataForm = ko.observable();

var stsTabForm = ko.observable("details");

var itemResource = ko.observableArray([]);

var listView = {}
listView.countryItems = ko.observableArray([]);
listView.systemItems = ko.observableArray([]);

listView.change = function() {
  var data = this.dataSource.view(),
  selected = $.map(this.select(), function(item) {
    return data[$(item).index()].ProductName;
  });

  console.log("data", data);
  console.log(
    this.select().map(function(item) {
      return data[$(item).index()];
    })
  );

  console.log("Selected: " + selected.length + " item(s), [" + selected.join(", ") + "]");
}

var stageSource = ko.observableArray([]);
var baseUriSource = ko.observableArray([]);
var domainSource = ko.observableArray([]);
var subDomainSource = ko.observableArray([]);

var resourceTemplate = {
  id: "",
  code: "",
  name: "",
  version: "",
  tag: [],
  description: "",
  raml: 0,
  stage: null,
  baseuri: null,
  domain: null,
  subdomain: null,
  parent: "",
  prodcoveragecountry :[],
  prodcoveragesystem :[],
  testcoveragecountry :[],
  testcoveragesystem :[]
}

var listVersion = ko.observableArray([]);
var selectedVersion = ko.observable(null)

var selectedResource = ko.computed(function(){
  return itemResource().find(function(v){
    return v.version() == selectedVersion()
  })
})

var selectedTag = ko.computed({
  read: function () {
    if(selectedResource()){
      return selectedResource().tag().join(", ");
    } else {
      return "";
    }
  },
  write: function(value) {
    if(selectedResource()){
      selectedResource().tag(value.split(", "));
    } else {
      return [];
    }
  }
})

stsForm.subscribe(function(){
  syncform();
});

listVersion.subscribe(function(){
  selectedVersion(listVersion()[listVersion().length - 1])
})

selectedResource.subscribe(function(v){
  if(v){
    ramldesigner.ramlAddr(filesUrl + "/" + v.parent() + "/" + (v._id ? v._id() : v.id()));

    if( ! v.baseuri()) v.baseuri(baseUriSource()[0]);
    if( ! v.domain()) v.domain(domainSource()[0]);
    if( ! v.subdomain()) v.subdomain(subDomainSource()[0]);

    v.uri = ko.computed(function(){
      return "https://" + v.baseuri() + (v.baseuri() ? ((v.domain() ? "/" : "") + v.domain() + (v.domain() ? ((v.subdomain() ? "/" : "") + v.subdomain()) : "")) : "");
    })

    v.uri.subscribe(doneTyping);
    doneTyping(v.uri());

    if($("li.active #menu-raml").length > 0){
      refreshEditor();
    }
  }
})

function onChangeProdCountry(e){
  // var selectedData = $('#idProdCoverageCountry').data('kendoListView').select()
  _.each($(".k-state-selected").find("p"), function(key, val){
      console.log(key.innerHTML)
  })
}

function saveResource(){
  pageProcessing(true);
  var param = ko.mapping.toJS(itemResource);

  param.forEach(function(v){
    v.id = v._id
  });

  ajaxPost("/resource/save", param, function(data){
    backToGrid();
  });
}

function deleteResource(){
  pageProcessing(true);
  var param = {
    parent: selectedResource().parent()
  }

  ajaxPost("/resource/delete", param, function(data){
    backToGrid();
  })
}

var goBack = function(){
  pageProcessing(true);
  if(stsDataForm() == "new"){
    deleteResource();
  } else {
    backToGrid();
  }
}

var handshake = function(){
  var param = ko.mapping.toJS(itemResource);

  param.forEach(function(v){
    v.id = v._id
  });

  ajaxPost("/resource/save", param, function(res){
    param = {};
    param.id = res.data[0].Parent

    ajaxPost("/resource/getlistresourcebycode", param, function(res){
      var protocol = 'https://';

      res.Data.forEach(function(v, i){
        res.Data[i].baseuri = v.uri.replace(protocol, "").split("/")[0];
        res.Data[i].domain = v.uri.replace(protocol, "").split("/")[1];
        res.Data[i].subdomain = v.uri.replace(protocol, "").split("/")[2];
      });

      var data = ko.mapping.fromJS(res.Data)
      itemResource(data())

      var ddVersion = $("#dd-version").data("kendoDropDownList");
      var selectedIndex = ddVersion.dataSource.data().indexOf(selectedVersion());
      ddVersion.select(selectedIndex);

      pageProcessing(false);
    })
  });
}

function callListVerison(param, callback){
  ajaxPost("/resource/getlistversionbycode", param, function(data){
    if(data.Data.length == 0){
      createNewVersion("1.0");
    } else {
      listVersion(data.Data);
    }

    stsForm('form');

    if(callback) callback();
  })
}

function newData(){
  pageProcessing(true);

  stsDataForm('new')
  stsTabForm("details")
  $("#txtCode").prop('disabled', false);

  callListVerison(resourceTemplate, handshake);
}

function showDetails(e) {
  e.preventDefault();
  pageProcessing(true);

  var dataItem = this.dataItem($(e.currentTarget).closest("tr"));
  dataItem.id = dataItem._id;

  var param = ko.mapping.fromJS(dataItem)

  callListVerison(param, function(){
    ajaxPost("/resource/getlistresourcebycode", param, function(res){
      var protocol = 'https://';

      res.Data.forEach(function(v, i){
        res.Data[i].baseuri = v.uri.replace(protocol, "").split("/")[0];
        res.Data[i].domain = v.uri.replace(protocol, "").split("/")[1];
        res.Data[i].subdomain = v.uri.replace(protocol, "").split("/")[2];
      });

      var data = ko.mapping.fromJS(res.Data)
      itemResource(data())

      var ddVersion = $("#dd-version").data("kendoDropDownList");
      var selectedIndex = ddVersion.dataSource.data().indexOf(selectedVersion());
      ddVersion.select(selectedIndex);

      pageProcessing(false);
    })
  });

  stsDataForm('exists');
  stsTabForm("details");
  $("#txtCode").prop('disabled', true);
}

var refreshGrid = function(){
  $('#grid').data('kendoGrid').dataSource.read();
  $('#grid').data('kendoGrid').refresh();
}

var backToGrid = function(){
  listVersion([]);
  itemResource([]);

  stsForm('grid');
  refreshGrid();
}

var createNewVersion = function(versionVal){
  var lastResource = itemResource()[itemResource().length - 1];

  var lastRes = ko.mapping.toJS(lastResource ? lastResource : resourceTemplate);
  lastRes._id = "";
  lastRes.parent = "";

  var res = ko.mapping.fromJS(lastRes);
  res.version(versionVal);

  listVersion.push(versionVal);
  itemResource.push(res);
}

function callGrid(){
  pageProcessing(true);

  $("#grid").kendoGrid({
    dataSource: {
      transport: {
        read: function(o){
          ajaxPost("/resource/getallresource", {}, function (res) {
            pageProcessing(false);
            o.success(res.Data)
          })
        }
      }
    },
    height: 610,
    sortable: true,
    pageable: {
      refresh: true,
      pageSizes: true,
      buttonCount: 5
    },
    columns: [{
      width: 50,
      title: "#",
      template: "<span class='row-number'></span>"
    }, {
      field: "uri",
      title: "URI"
    }, {
      field: "raml",
      title: "RAML",
      template: "#= raml # files"
    }, {
      field: "version",
      title: "Current Version",
    },
    // {
    //   field: "stage",
    //   title: "Stages",
    //   template: function(v){
    //     var stage = stageSource().find(function(s){
    //       return s._id == v.stage
    //     })
    //
    //     return stage ? stage.name : ""
    //   }
    // },
    {
      command: { text: "Edit", click: showDetails }, title: "Action", width: "85px"
    }],
    dataBound: function () {
      var rows = this.items();
      $(rows).each(function () {
        var index = $(this).index() + 1;
        var rowLabel = $(this).find(".row-number");
        $(rowLabel).html(index);
      });
    }
  });
}

var syncform = function(){
  ramldesigner.buttondetailvis(false);
}

var tabdetailclick = function(){
  stsTabForm("details");
  syncform();
}

var tablogclick = function(){
  stsTabForm("log");
  syncform();
}

$(document).ready(function () {
  callGrid();

  ajaxPost("/masterstage/getallstage", { }, function(res){
  	stageSource(res.Data);
  });

  ajaxPost("/masteruri/getdata", { }, function(res){
    var baseuri = res.data.find(function(v){ return v.Name == "BASE_URI" });
    if(baseuri) baseUriSource(baseuri.Options);

    var domain = res.data.find(function(v){ return v.Name == "DOMAIN" });
    if(domain) domainSource(domain.Options);

    var subdomain = res.data.find(function(v){ return v.Name == "SUB_DOMAIN" });
    if(subdomain) subDomainSource(subdomain.Options);
  });

  ajaxPost("/mastercountry/getallcountry", { }, function(res){
    listView.countryItems(res.Data);
  });
  ajaxPost("/mastersystem/getallsystem", { }, function(res){
    listView.systemItems(res.Data);
  });
  listView.template = kendo.template('<div class="list-option"><p>#= name #</p></div>');
});

//user is "finished typing," do something
function doneTyping(val) {
  if(val){
    val = val.toLowerCase();

    var protocol = 'https://';

    var domain = ((function(addr){
      return addr.replace(protocol, "").split("/")[0];
    })(val));

    var checkAlphanumeric = ((function(d){
      var ret = true;
      var exp = /^([0-9]|[a-z])+([0-9a-z]+)$/i;

      d.split(".").forEach(function(v){
        if( ! v.match(exp)){
          ret = false
        }
      });

      return ret;
    })(domain));

    var checkHttps = (val.substring(0, 8) == protocol);
    var checkDomain = (domain.endsWith(".sc.com"));
    var checkLowercase = (val == val.toLowerCase());
    var checkNoUnderscore = (val.indexOf("_") == -1);
    var checkForwardSlash = ( ! val.endsWith("/"));
    var checkJson = (val.indexOf(".json") == -1)
    var checkXml = (val.indexOf(".xml") == -1)

    var checkVersion = ((function(d){
      var ret = false;
      var exp = /\bv\d{1,2}\b/g;

      d.split("/").forEach(function(v){
        if(v.match(exp)){
          ret = true
        }
      });

      return ret;
    })(val));

    if (checkAlphanumeric && checkHttps && checkDomain && checkLowercase && checkNoUnderscore && checkForwardSlash && checkJson && checkXml && checkVersion){
      stsExistUri(true);
    } else {
      stsExistUri(false);
    }

    return stsExistUri()

    // $.ajax({
    //   url: val,
    //   //data: myData,
    //   type: 'GET',
    //   crossDomain: true,
    //   dataType: 'jsonp',
    //   complete: function(res) {
    //     var resStat = res.status
    //     if (resStat == 200){
    //       stsExistUri(true);
    //     } else {
    //       stsExistUri(false);
    //     }
    //   },
    // });
  }
}

$('#modal-version').on('hidden.bs.modal', function () {
  $("#grid-version").html("");
});

$('#modal-version').on('shown.bs.modal', function (e) {
  var mapVersion = function(){
    return listVersion().map(function(v){
      return { version: v, id: v };
    })
  }

  var refreshGridVersion = function(){
    $('#grid-version').data('kendoGrid').dataSource.read();
    $('#grid-version').data('kendoGrid').refresh();
  }

  $("#grid-version").kendoGrid({
    dataSource: {
      transport: {
        read: function(o){
          o.success(mapVersion());
        },
        update: function(o){
          listVersion.replace(o.data.id, o.data.version);

          var found = itemResource().find(function(v){
            return v.version() == o.data.id
          })
          found.version(o.data.version);

          o.success(mapVersion());
          refreshGridVersion();
        },
        create: function(o){
          pageProcessing(true);
          createNewVersion(o.data.version);

          o.success(mapVersion());
          refreshGridVersion();

          handshake();
        }
      },
      schema: {
        model: {
          id: "id",
          fields: {
            version: { validation: { required: true } },
          }
        }
      },
    },
    toolbar: ["create"],
    editable: "inline",
    columns: [{
      width: 50,
      title: "#",
      template: "<span class='row-number'></span>"
    }, {
      field: "version",
      title: "Version"
    },
    { command: ["edit"], title: "&nbsp;", width: "85px" }],
    dataBound: function () {
      var rows = this.items();
      $(rows).each(function () {
        var index = $(this).index() + 1;
        var rowLabel = $(this).find(".row-number");
        $(rowLabel).html(index);
      });
    }
  })
})
