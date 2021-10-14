import draftJS from 'draft-js';

let draftState = process.argv.slice(2).join(" ");
let jsonDraftState = JSON.parse(draftState);
let contentState = draftJS.ContentState.createFromBlockArray(jsonDraftState.blocks, jsonDraftState.entityMap)
let editorState = draftJS.EditorState.createWithContent(contentState);

let rawState = draftJS.convertToRaw(editorState.getCurrentContent());
console.log(JSON.stringify(rawState));