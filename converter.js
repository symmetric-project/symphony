import { markdownToDraft } from 'markdown-draft-js';

let rawState = process.argv.slice(2).toString();
console.log('Raw state: ' + rawState);

var markdownString = markdownToDraft(rawState);

console.log(markdownString);