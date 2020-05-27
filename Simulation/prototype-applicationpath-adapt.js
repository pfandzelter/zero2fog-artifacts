'use strict';


// load fogexplorer node module
require("./FogExplorer/node/nodeGlobalizer");
const evaluator = require("./FogExplorer/sharedLogic/evaluator");
const fs = require('fs');

// create different mappings

let adapt = ["sensor", "pack_cntrl", "wgw", "fctry_datacenter"];


let mappingOptions = [];

adapt.forEach((c, i) => {
  mappingOptions.push({
    c: c,
  });
});


// prepare csv
const filename = "results-applicationpath-adapt.csv";
fs.writeFileSync(filename, "camera_type,wgw_type,fctry_datacenter_type,office_datacenter_type,cloud_type,prod_cntrl_type,adapt_map,proc_time,proc_cost,tx_time,tx_cost,full_time,full_cost,existsCoherentFlow\n");

// create different options

let camera = [{
  availableMemory: 1,
  memoryPrice: 500,
  performanceIndicator: 0.001
}, {
  availableMemory: 10,
  memoryPrice: 2500,
  performanceIndicator: 0.05
}];

let prod_cntrl = [{
  availableMemory: 10,
  memoryPrice: 250,
  performanceIndicator: 0.05
}, {
  availableMemory: 100,
  memoryPrice: 1000,
  performanceIndicator: 0.5
}];

let wgw = [{
  availableMemory: 100,
  memoryPrice: 250,
  performanceIndicator: 0.1
}, {
  availableMemory: 250,
  memoryPrice: 500,
  performanceIndicator: 0.5
}];

let fctry_datacenter = [{
  availableMemory: 5000,
  memoryPrice: 5,
  performanceIndicator: 0.5
}, {
  availableMemory: 10000,
  memoryPrice: 10,
  performanceIndicator: 1
}, {
  availableMemory: 20000,
  memoryPrice: 25,
  performanceIndicator: 5
}];

let office_datacenter = [{
  availableMemory: 50000,
  memoryPrice: 5,
  performanceIndicator: 5
}, {
  availableMemory: 100000,
  memoryPrice: 10,
  performanceIndicator: 10
}, {
  availableMemory: 200000,
  memoryPrice: 20,
  performanceIndicator: 20
}];

let cloud = [{
  availableMemory: 50000,
  memoryPrice: 1,
  performanceIndicator: 10
}, {
  availableMemory: 100000,
  memoryPrice: 2,
  performanceIndicator: 25
}, {
  availableMemory: 200000,
  memoryPrice: 5,
  performanceIndicator: 50
}];

let infrastructureOptions = [];

camera.forEach((c, ci) => {
  prod_cntrl.forEach((p, pi) => {
    wgw.forEach((w, wi) => {
      fctry_datacenter.forEach((f, fi) => {
        office_datacenter.forEach((o, oi) => {
          cloud.forEach((cl, cli) => {
            infrastructureOptions.push({
              c: c,
              w: w,
              f: f,
              o: o,
              cl: cl,
              p: p,
              ci: ci,
              wi: wi,
              fi: fi,
              oi: oi,
              cli: cli,
              pi: pi
            });
          });
        });
      });
    });
  });
});

// calculate metrics for each one

let counter = 0;

infrastructureOptions.forEach((infrastructure, i) => {
  // create infrastructure model as object based on based model
  // i hate js
  let thismodel = JSON.parse(JSON.stringify(require("./model-adapt.json")));

  thismodel.nodes.push({
    id: "camera",
    label: "Camera",
    baseProperties: infrastructure.c
  });

  thismodel.nodes.push({
    id: "prod_cntrl",
    label: "Production Controller",
    baseProperties: infrastructure.p
  });

  thismodel.nodes.push({
    id: "wgw",
    label: "Wireless Gateway",
    baseProperties: infrastructure.w
  });

  thismodel.nodes.push({
    id: "fctry_datacenter",
    label: "Factory Data Center",
    baseProperties: infrastructure.f
  });

  thismodel.nodes.push({
    id: "office_datacenter",
    label: "Central Office Data Center",
    baseProperties: infrastructure.o
  });

  thismodel.nodes.push({
    id: "cloud",
    label: "Cloud",
    baseProperties: infrastructure.cl
  });

  //load model
  evaluator.loadModels(thismodel);

  mappingOptions.forEach((mapping, i) => {
    counter++
    console.log(counter)
    // add mapping
    evaluator.setCurrentlySelectedModuleToModuleWithId("adapt");
    evaluator.assignCurrentlySelectedModuleToNodeWithId(mapping.c);

    // calculate metrics
    let existsCoherentFlow = evaluator.calculateCoherentFlow();
    console.log(existsCoherentFlow)
    let procTime = round(evaluator.getTotalProcessingTime(), 5);
    let procCost = round(evaluator.getTotalProcessingCost(), 5);
    let txTime = round(evaluator.getTotalTransmissionTime(), 5);
    let txCost = round(evaluator.getTotalTransmissionCost(), 5);

    // write to csv
    let resultsString = "";
    resultsString += infrastructure.ci.toString();
    resultsString += ",";
    resultsString += infrastructure.wi.toString();
    resultsString += ",";
    resultsString += infrastructure.fi.toString();
    resultsString += ",";
    resultsString += infrastructure.oi.toString();
    resultsString += ",";
    resultsString += infrastructure.cli.toString();
    resultsString += ",";
    resultsString += infrastructure.pi.toString();
    resultsString += ",";
    resultsString += mapping.c;
    resultsString += ",";
    resultsString += procTime;
    resultsString += ",";
    resultsString += procCost;
    resultsString += ",";
    resultsString += txTime;
    resultsString += ",";
    resultsString += txCost;
    resultsString += ",";
    resultsString += procTime + txTime;
    resultsString += ",";
    resultsString += procCost + txCost;
    resultsString += ",";
    resultsString += existsCoherentFlow;
    resultsString += "\n";

    fs.appendFileSync(filename, resultsString);

    if (existsCoherentFlow && (procTime + txTime != Number.POSITIVE_INFINITY) && (procCost + txCost != Number.POSITIVE_INFINITY)) {
      let name = "adapt-";
      name += infrastructure.ci.toString();
      name += "-";
      name += infrastructure.wi.toString();
      name += "-";
      name += infrastructure.fi.toString();
      name += "-";
      name += infrastructure.oi.toString();
      name += "-";
      name += infrastructure.cli.toString();
      name += "-";
      name += infrastructure.pi.toString();
      name += "-";
      name += mapping.c;

      fs.writeFileSync("./results-applicationpath/" + name + ".json", evaluator.getAllDataRepresentation());
    }
  });
});