;"use strict";
class Starter {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  mover(e) {
    const ta = e.currentTarget;
    const txt = ta.textContent;
    document.querySelectorAll('td.hita').forEach(el => {
      el.classList.remove('se');
      const se = el.textContent;
      if (se !== txt) return;
      el.classList.add('se');
    });
  }

  hiclick(e) {
    const ta = e.currentTarget;
    const buf = ta.textContent;
    window.navigator.clipboard.writeText(buf).then(() => {
      console.log(buf + ' copied');
    });
  }

  handler() {
    const refLink = document.getElementById('ref-ttl');

    if (!refLink) return;

    document.querySelectorAll('td.hita').forEach(el => {
      el.addEventListener('mouseover', this.mover.bind(this));
      el.addEventListener('click', this.hiclick.bind(this));
    });
    refLink.addEventListener('click', () => {
      window.location.reload();
    });
  }
}

new Starter();
