import { convertToRaw } from 'draft-js';

let draftState = process.argv.slice(2).join(" ");
/* console.log('Draft state: ' + draftState); */
let rawState = convertToRaw(draftState);
console.log(JSON.stringify(rawState));