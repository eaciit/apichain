"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var fs = require("fs");
var def = require("raml-definition-system");
var path = require("path");
var mkdirp = require("mkdirp");
var docGen = require("./docGen");
var universe10 = def.getUniverse("RAML10");
var universe08 = def.getUniverse("RAML08");
var apiType10 = universe10.type("Api");
var apiType08 = universe08.type("Api");
var documentationFolder = "../../../documentation";
var docTablesFile10 = path.join(__dirname, documentationFolder + "/RAML10Classes.html").toString();
var docTablesFile08 = path.join(__dirname, documentationFolder + "/RAML08Classes.html").toString();
mkdirp.sync(path.dirname(docTablesFile10));
fs.writeFileSync(docTablesFile10, docGen.def2Doc(apiType10));
fs.writeFileSync(docTablesFile08, docGen.def2Doc(apiType08));
//# sourceMappingURL=docGenScript.js.map