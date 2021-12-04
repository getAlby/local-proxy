document.getElementById("address").focus();

document.getElementById('address').value = localStorage.getItem('address');
document.getElementById('cert').value = localStorage.getItem('cert');
document.getElementById('port').value = localStorage.getItem('port') || '8181';
document.getElementById('secret').value = localStorage.getItem('secret') || Math.random().toString(20).substr(2, 10);

document.getElementById("send").addEventListener("click", function (e) {
  e.preventDefault();
  if (!window.proxyRunning) {
    window.StartProxy();
    e.target.innerText = "Stop";
    e.target.style.display = "none";
  } else {
    //window.StopProxy();
    //e.target.innerText = "Start";
  }
});

// TODO: does not work in go
window.StartProxy = function () {
  document.getElementById("result").innerText = "starting...";
  window.proxyRunning = true;
  const address = document.getElementById('address').value;
  const cert = document.getElementById('cert').value;
  const port = document.getElementById('port').value;
  const secret = document.getElementById('secret').value;
  localStorage.setItem('address', address);
  localStorage.setItem('cert', cert);
  localStorage.setItem('port', port);
  localStorage.setItem('secret', secret);

  window.go.main.App.StartProxy(address, secret, cert, port).then((result) => {
    document.getElementById("result").innerText = result;
    document.getElementById("send").innerText = "Stop";
  });
};

window.StopProxy = function () {
  window.proxyRunning = false;
  window.go.main.App.StopProxy().then((result) => {
    document.getElementById("result").innerText = result;
    document.getElementById("send").innerText = "Start";
  });
}

