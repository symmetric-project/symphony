import { markdownToDraft } from 'markdown-draft-js';
let markdown = process.argv.slice(2).join(" ");
let draftState = markdownToDraft(markdown);

if (draftState.blocks.length === 1 && draftState.blocks[0].text === "") {
} else {
    let draftState = markdownToDraft(markdown);
    console.log(JSON.stringify(draftState));
}