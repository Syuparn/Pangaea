const go = new Go();

// run wasm
fetch("./main.wasm").then(response => 
  response.arrayBuffer()
).then(bytes =>
  WebAssembly.instantiate(bytes, go.importObject)
).then(obj => {
  go.run(obj.instance);
  // HACK: replace "now loading..." with sample code
  document.getElementById('source').value = `"Hello, world!".p`;
});
