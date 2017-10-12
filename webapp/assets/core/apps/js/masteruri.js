master = {};
master.record = ko.observable({
  _id: ko.observable(""),
  StringOptions : ko.observable(""),
  Name: ko.observable("")
})

master.new = function(){
  master.record(ko.mapping.fromJS({
    _id: "",
    StringOptions : "",
    Name: ""
  }))
}

master.getData = function(){
  ajaxPost("getdata", {}, function(res) {
    if (res.success == true){
      master.generateGrid(res.data)
    }else{
      return swal("Error!", res.Message, "error");
    }
  });
}

master.saveData = function() {
  ajaxPost("savedata", master.record(), function(data) {
    if (data.success == true){
      master.record(ko.mapping.fromJS({
        _id: "",
        StringOptions : "",
        Name: ""
      }))

      master.getData()
    }else{
      return swal("Error!", data.Message, "error");
    }
  });
}

master.deleteData = function(e) {
  ajaxPost("deletedata", {_id: e}, function (data) {
    if (data.success == true){
      master.record(ko.mapping.fromJS({
        _id: "",
        StringOptions : "",
        Name: ""
      }))
      master.getData()
    }else{
      return swal("Error!", data.Message, "error");
    }
  });
}

master.editConfig = function(id) {
  var grid = $("#grid-master").getKendoGrid()._data;
  var rec = grid.find(x => x._id == id);
  if (rec != undefined) {
    master.record(ko.mapping.fromJS(rec))
  }
}

master.generateGrid = function(data) {
  $("#grid-master").html("");
  $("#grid-master").kendoGrid({
    dataSource: {
      data: data
    },
    // height: 550,
    columns: [{
        template: function(dataItem){
          return '<button class="btn btn-sm btn-warning btn-fill" onclick="master.editConfig(\''+ dataItem._id +'\')"><span class="glyphicon glyphicon-edit"></span></button>&nbsp;&nbsp;<button class="btn btn-sm btn-danger btn-fill" onclick="master.deleteData(\''+ dataItem._id +'\')"><span class="glyphicon glyphicon-trash"></span></button>'
        },
        width: 100
      },
      { field: "Name", title: "Name" },
      { field: "Options", title: "Options", template: function(e){ return e.Options.join(", ") } },
    ]
  })
}

$(document).ready(function(){
  master.getData()
})
