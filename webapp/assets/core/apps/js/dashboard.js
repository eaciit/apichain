var dashboard = {}

dashboard.viewTypeChoices = ko.observableArray(["Business View", "Developer View", "User Journey", "Option"]);
dashboard.viewTypeVal = ko.observable();
dashboard.cb = ko.observableArray([]);

dashboard.busGraphVis = ko.computed(function(){ return dashboard.viewTypeVal() == "Business View" });
dashboard.devGraphVis = ko.computed(function(){ return dashboard.viewTypeVal() == "Developer View" });
dashboard.journeyGraphVis = ko.computed(function(){ return dashboard.viewTypeVal() == "User Journey" });

dashboard.renderDeveloperGraph = function(){
  var nodes = new vis.DataSet([
    {id: 1, label: 'Cyber Security', shape: 'circle'},
    {id: 2, label: 'Digital Product', shape: 'circle'},
    {id: 3, label: 'Omni', shape: 'circle'},
    {id: 4, label: 'Host 2 host', shape: 'circle'},
    {id: 5, label: ' ', shape: 'circle'},
    {id: 6, label: ' ', shape: 'circle'},
    {id: 7, label: ' ', shape: 'circle'},
    {id: 8, label: ' ', shape: 'circle'},
    {id: 9, label: ' ', shape: 'circle'},
    {id: 10, label: ' ', shape: 'circle'},
    {id: 11, label: ' ', shape: 'circle'},
    {id: 12, label: ' ', shape: 'circle'},
    {id: 13, label: ' ', shape: 'circle'},
    {id: 14, label: ' ', shape: 'circle'},
  ]);

  var edges = new vis.DataSet([
    {from: 1, to: 5, color: { color: 'black' }},
    {from: 1, to: 6, color: { color: 'black' }},
    {from: 1, to: 7, color: { color: 'black' }},
    {from: 1, to: 8, color: { color: 'black' }},
    {from: 1, to: 9, color: { color: 'black' }},
    {from: 1, to: 3, color: { color: 'black' }},
    {from: 2, to: 3, color: { color: 'black' }},
    {from: 2, to: 5, color: { color: 'black' }},
    {from: 5, to: 3, color: { color: 'black' }},
    {from: 9, to: 10, color: { color: 'black' }},
    {from: 10, to: 3, color: { color: 'black' }},
    {from: 3, to: 4, color: { color: 'black' }},
    {from: 3, to: 11, color: { color: 'black' }},
    {from: 3, to: 12, color: { color: 'black' }},
    {from: 2, to: 12, color: { color: 'black' }},
    {from: 4, to: 13, color: { color: 'black' }},
    {from: 13, to: 14, color: { color: 'black' }},
    {from: 14, to: 3, color: { color: 'black' }},
  ]);

  var container = document.getElementById('graph-developer');

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
}

dashboard.renderBusinessGraph = function(){
  var simple_chart_config = {
    chart: {
      container: "#graph-business",
      nodeAlign: "BOTTOM",
      rootOrientation: "WEST",
      connectors: {
        style: {
					'stroke': '#A06C31',
          'stroke-width': 2
				}
			},
    },
    nodeStructure: {
      text: { name: "CIB" },
      HTMLclass: "bunder",
      collapsable: true,
      children: [
        {
          text: { name: "Omni Channel" },
          HTMLclass: "kotak",
          collapsable: true,
          children: [
            {
              text: { name: "Self Service" },
              drawLineThrough: true,
              collapsable: true,
              children: [
                {
                  text: { name: "Interactive Web" },
                  drawLineThrough: true,
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 1" },
                      drawLineThrough: true,
                    },
                    {
                      text: { name: "Resource 2" },
                      drawLineThrough: true,
                    }
                  ]
                },
                {
                  text: { name: "Host to Host" },
                  drawLineThrough: true,
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 3" },
                      drawLineThrough: true,
                    }
                  ]
                },
                {
                  text: { name: "API" },
                  drawLineThrough: true,
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 1" },
                      drawLineThrough: true,
                    },
                    {
                      text: { name: "Resource 4" },
                      drawLineThrough: true,
                    }
                  ]
                }
              ]
            },
            {
              text: { name: "Assisted Service" },
              drawLineThrough: true,
              collapsable: true,
              children: [
                {
                  text: { name: "Assisted Trx Initation" },
                  drawLineThrough: true,
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 5" },
                      drawLineThrough: true,
                    }
                  ]
                },
                {
                  text: { name: "Supporting Document Validation" },
                  drawLineThrough: true,
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 1" },
                      drawLineThrough: true,
                    },
                    {
                      text: { name: "Resource 5" },
                      drawLineThrough: true,
                    },
                  ]
                },
              ]
            },
            {
              text: { name: "Industry & Country Networks" },
              drawLineThrough: true,
              collapsable: true,
              children: [
                {
                  text: { name: "S.W.I.F.T" },
                  drawLineThrough: true,
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 3" },
                      drawLineThrough: true,
                    }
                  ]
                }
              ]
            }
          ]
        },
        {
          text: { name: "Cyber Security" },
          HTMLclass: "kotak",
          collapsable: true,
          children: [
            {
              text: { name: "Identity" },
              drawLineThrough: true,
              collapsable: true,
              children: [
                {
                  text: { name: "Unified User Model" },
                  drawLineThrough: true,
                },
                {
                  text: { name: "Tokens" },
                  drawLineThrough: true,
                },
                {
                  text: { name: "Crypto" },
                  drawLineThrough: true,
                },
                {
                  text: { name: "Biometrics" },
                  drawLineThrough: true,
                }
              ]
            },
            {
              text: { name: "Access" },
              drawLineThrough: true,
              collapsable: true,
              children: [
                {
                  text: { name: "Fraud Monitoring" },
                  drawLineThrough: true,
                },
                {
                  text: { name: "Perimeter Defense" },
                  drawLineThrough: true,
                },
              ]
            },
            {
              text: { name: "Information" },
              drawLineThrough: true,
              collapsable: true,
              children: [
                {
                  text: { name: "Entitlements" },
                  drawLineThrough: true,
                },
                {
                  text: { name: "Digital Signatures" },
                  drawLineThrough: true,
                },
              ]
            }
          ]
        },
        {
          text: { name: "Digital Product" },
          HTMLclass: "kotak",
        },
        {
          text: { name: "Data" },
          HTMLclass: "kotak",
        }
      ]
    }
  };

  new Treant( simple_chart_config );
}

dashboard.renderUserJourneyGraph = function(){
  var simple_chart_config = {
    chart: {
      container: "#graph-journey",
      nodeAlign: "BOTTOM",
      rootOrientation: "WEST",
      connectors: {
        style: {
          'stroke': '#808080',
          'stroke-width': 2,
        }
      },
    },
    nodeStructure: {
      text: { name: "Trade Monitor System" },
      HTMLclass: "kotakbgwhite",
      collapsable: false,
      connectors: {
        style: {
          'stroke': '#808080',
          'stroke-width': 2,
          'arrow-end': 'butt',
        }
      },
      children: [

        {
          text: { name: "Login" },
          HTMLclass: "kotakbgblue",
          collapsable: true,
          children :[{
          text: { name: "Dashboard" },
          HTMLclass: "kotakbgblue",
          collapsable: true,
          children: [
            {
              text: { name: "Portofolio" },
              HTMLclass: "kotakbgblue",
              collapsable: true,
              children: [
                {
                  text: { name: "Trade" },
                  HTMLclass: "kotakbgblue",
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 1" },
                      drawLineThrough: true,
                    },
                    {
                      text: { name: "Resource 2" },
                      drawLineThrough: true,
                    }
                  ]
                },
                {
                  text: { name: "Resource 3" },
                  drawLineThrough: true,
                  collapsable: true,
                },
                {
                  text: { name: "Resource 4" },
                  drawLineThrough: true,
                  collapsable: true,

                }
              ]
            },
            {
              text: { name: "Catalogue Browse" },
             HTMLclass: "kotakbgblue",
              collapsable: true,
              children: [

                {
                  text: { name: "Trade" },
                  HTMLclass: "kotakbgblue",
                  collapsable: true,
                  children: [
                    {
                      text: { name: "Resource 1" },
                      drawLineThrough: true,
                    },
                    {
                      text: { name: "Resource 2" },
                      drawLineThrough: true,
                    },
                  ]
                },
              ]
            },
            {
              text: { name: "Trade" },
              HTMLclass: "kotakbgblue",
            },
            {
              text: { name: "Deposit" },
               HTMLclass: "kotakbgblue",
            },
            {
              text: { name: "Withdraw" },
               HTMLclass: "kotakbgblue",

            }
          ]
        }]

        }
      ]
    }
  };

  new Treant( simple_chart_config );
}

dashboard.renderGraph = function(){
  switch (dashboard.viewTypeVal()) {
    case "Developer View":
    dashboard.renderDeveloperGraph();
    break;
    case "Business View":
    dashboard.renderBusinessGraph();
    break;
    case "User Journey":
    dashboard.renderUserJourneyGraph();
    break;
  }
}

dashboard.viewTypeChange = function(asfd){
  setTimeout(dashboard.renderGraph, 0);
};

$(document).ready(function(){
  dashboard.viewTypeVal(dashboard.viewTypeChoices()[1]);
  dashboard.renderGraph();
});
