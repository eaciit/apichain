var fs = require('fs');

exports.require = function(path,names_to_export) {
  filedata = fs.readFileSync(path,'utf8');
  eval(filedata);
  exported_obj = {};
  for (i in names_to_export) {
    to_eval = 'exported_obj[names_to_export[i]] = '
    + names_to_export[i] + ';'
    eval(to_eval);
  }
  return exported_obj;
}
