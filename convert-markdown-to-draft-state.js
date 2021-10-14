import { markdownToDraft } from 'markdown-draft-js';

let markdown = process.argv.slice(2).join(" ");
/* console.log('Markdown: ' + markdown); */
let draftState = markdownToDraft(markdown);
console.log(JSON.stringify(draftState));