class Funcs {
  ready(hin) {
    if (document.readyState != "loading") {
      hin();
      return
    }

    document.addEventListener("DOMContentLoaded", hin);
  }

  parent(el, match) {
    let _el = el.parentNode;

    while (_el && typeof _el.matches === 'function') {
      if (_el.matches(match)) return _el;

      _el = _el.parentNode;
    }

    return null;
  }

  position(el) {
    const re = el.getBoundingClientRect();
    const cs = getComputedStyle(el);

    return {
      "top": re.top - parseInt(cs.marginTop, 10),
      "left": re.left - parseInt(cs.marginLeft, 10)
    };
  }

  trigger(el, ev_string) {
    const ev = new CustomEvent(ev_string);
    el.dispatchEvent(ev);
  }

  once(el, ev_name) {
    if (el.data && el.data[ev_name]) return true;

    if (!el.data) {
      el.data = {};
    }

    if (!el.data[ev_name]) {
      el.data[ev_name] = true;
    }

    return false;
  }
}
