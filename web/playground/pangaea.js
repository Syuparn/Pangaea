const go = new Go();

// run wasm
fetch("./main.wasm").then(response => 
  response.arrayBuffer()
).then(bytes =>
  WebAssembly.instantiate(bytes, go.importObject)
).then(obj => {
  go.run(obj.instance);
  // HACK: replace "now loading..." with sample code
  document.getElementById('source').value = initializeSourceCode();
});


function initializeSourceCode() {
  const fragment = trimURIFragment();
  if (fragment === '') {
    return `"Hello, world!".p`;
  }
  return decodeURI(fragment);
}

function trimURIFragment() {
  // trim prefix "#"
  // NOTE: replace replaces only the first occurence
  return location.hash.replace('#', '');
}
