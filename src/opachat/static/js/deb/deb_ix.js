;"use strict";
$(() => {
  const loop = (list, tr, tb, pa) => {
    if (!list) return;

    list.forEach((ii) => {
      let tra = tr
        .replace(/#CAP#/g, ii.cap)
        .replace(/#PAA#/g, pa);

      tb.append($(tra));

      loop(ii.list, tr, tb, pa + 30);
    });
  };

  const load_deb = () => {
    let cs = document.getElementsByName("gorilla.csrf.Token")[0].value;
    let ws_main = $("#ws-main").eq(0);
    let tb = $('#tbl-info tbody').eq(0);
    let tr = cr_tag();

    tb.html('');

    let url = ws_main.attr('href');

    axios.post(url, null, {
      headers: { "X-CSRF-Token": cs }
    })
      .then(re => {
        loop(re.data.list, tr, tb, 0);
      })
      .catch(err => {
        console.log(err);
      });
  };

  const cr_tag = () => {
    let res = '<tr><td style="padding-left:#PAA#px">#CAP#</td></tr>'
    return res;
  };

  $('#info-ref').click(() => {
    load_deb();
  });

  load_deb();
});
