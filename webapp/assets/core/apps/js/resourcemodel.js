var resmodel = {
  id: "graph-resourcemodel",
  ramldata: ko.observableArray([])
}

resmodel.fetchResources = function(){
  var resources = [];

  resmodel.ramldata().forEach(function(v){
    if(v.resources){
      v.resources.forEach(function(w){
        resources.push(w.relativeUri);
      });
    }
  });

  return resources.filter(function(item, pos) {
    return resources.indexOf(item) == pos;
  });
}

resmodel.fetchSchemas = function(data){
  var foundSchemas = findValues(data, "schema").map(function(v){
  	return v[0];
  });

  return cleanArray(foundSchemas.filter(function(item, pos) {
    return foundSchemas.indexOf(item) == pos;
  }));
}

resmodel.fetchEdges = function(){
  var edges = [];

  resmodel.ramldata().forEach(function(v){
    if(v.resources){
      v.resources.forEach(function(resource){
        var foundSchemas = resmodel.fetchSchemas(resource);
        foundSchemas.forEach(function(schema){
          edges.push({ resource: resource.relativeUri, schema: schema })
        })
      });
    }
  });

  return edges.filter(function(item, pos) {
    return edges.indexOf(item) == pos;
  });
}

resmodel.renderGraph = function(){
  var resourcesDataset = resmodel.fetchResources().map(function(v){
    return { id: v, label: v, shape: 'circle' };
  });

  var schemasDataset = resmodel.fetchSchemas(resmodel.ramldata()).map(function(v){
    return { id: v, label: v, shape: 'square' };
  });

  var edgesDataset = resmodel.fetchEdges().map(function(v){
    return { from: v.resource, to: v.schema, color: { color: 'black' } };
  });

  var nodes = new vis.DataSet(resourcesDataset.concat(schemasDataset));

  var edges = new vis.DataSet(edgesDataset);

  var container = document.getElementById(resmodel.id);

  var data = {
    nodes: nodes,
    edges: edges
  };

  var options = {
    "interaction": {
      "hover": true
    },
    "edges": {
      "smooth": {
        "forceDirection": "none"
      }
    },
    "physics": {
      "forceAtlas2Based": {
        "avoidOverlap": 1
      },
      "minVelocity": 0.75,
      "solver": "forceAtlas2Based"
    }
  };

  var network = new vis.Network(container, data, options);

  network.on("hoverNode", function (params) {
    console.log('hoverNode Event:', params);
  });

  network.on("doubleClick", function (params) {
    console.log('doubleClick Event:', params);
  });

  network.on("stabilizationIterationsDone", function () {
    network.setOptions( { physics: false } );
  });
};

resmodel.render = function(){
  stsTabForm("resmodel");
  syncform();

  setTimeout(function(){
    $("#" + resmodel.id).html("");
    resmodel.ramldata([]);

    ajaxPost("/designer/getramllist", { parent: selectedResource().parent(), id: selectedResource()._id() }, function(res){
      var paths = res.data;

      paths.forEach(function(path){
        var fixedPath = window.location.origin + "/files" + path;
        pushRamlData(fixedPath, resmodel.ramldata);
      });

    	resmodel.renderGraph();
    });
  }, 200);
};
