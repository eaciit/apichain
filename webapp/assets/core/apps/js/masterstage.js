$(document).ready(function(){
  var refreshGrid = function(){
    $('#grid-main').data('kendoGrid').dataSource.read();
    $('#grid-main').data('kendoGrid').refresh();
  }

  $("#grid-main").kendoGrid({
    dataSource: {
      transport: {
        read: function(o){
          ajaxPost("/masterstage/getallstage", { }, function(res){
            o.success(res.Data);
          })
        },
        update: function(o){
          ajaxPost("/masterstage/save", { id: o.data._id, name: o.data.name }, function(res){
            o.success(res.Data);
            refreshGrid();
          });
        },
        create: function(o){
          ajaxPost("/masterstage/save", { id: o.data._id, name: o.data.name }, function(res){
            o.success(res.Data);
            refreshGrid();
          });
        },
        destroy: function(o){
          ajaxPost("/masterstage/delete", { id: o.data._id, name: o.data.name }, function(res){
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
