$(document).ready(function(){
  var refreshGrid = function(){
    $('#grid-main').data('kendoGrid').dataSource.read();
    $('#grid-main').data('kendoGrid').refresh();
  }

  $("#grid-main").kendoGrid({
    dataSource: {
      transport: {
        read: function(o){
          ajaxPost("/masterschema/getallschema", { }, function(res){
            console.log(res.Data[0].jsonstring.toString())
            o.success(res.Data);
          })
        },
        update: function(o){
          ajaxPost("/masterschema/save", { id: o.data._id, name: o.data.name, jsonstring: o.data.jsonstring }, function(res){
            o.success(res.Data);
            refreshGrid();
          });
        },
        create: function(o){
          ajaxPost("/masterschema/save", { id: o.data._id, name: o.data.name, jsonstring: o.data.jsonstring }, function(res){
            o.success(res.Data);
            refreshGrid();
          });
        },
        destroy: function(o){
          ajaxPost("/masterschema/delete", { id: o.data._id, name: o.data.name, jsonstring: o.data.jsonstring }, function(res){
            o.success(res.Data);
            refreshGrid();
          });
        },
      },
      schema: {
        model: {
          id: "_id",
          fields: {
            name: { validation: { required: true } },
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
      field: "name",
      title: "Name"
    }, {
      field: "jsonstring",
      title: "JSON"
    }, {
      command: ["edit", "destroy"], title: "&nbsp;", width: "175px"
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
})
