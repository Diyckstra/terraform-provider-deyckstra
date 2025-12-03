(function () {
  const CONFIG = {
    autoBySeparator: false, // переносить "заголовок" по завершающему символу (если потребуется)
    separatorRe: /[:：—–-]\s*$/,
    // Ключевые слова для автоматического отступа в note-блоках
    gapKeywords: ['note', 'warning', 'info', 'tip', 'important'],
    addGapAfterKeywordStrongInNotes: true
  };

  function processParagraphs() {
    document.querySelectorAll('p').forEach(function (paragraph) {
      if (paragraph.dataset.boldjsProcessed === '1') return;

      let html = paragraph.innerHTML;
      let iconHTML = '';
      let noteColor = ''; // 'yellow' | 'blue' | 'red'

      // Маркеры заметок в начале параграфа: ~>, ->, !> (и &gt;)
      const markerMatch = html.match(/^\s*(~(?:&gt;|>)|-(?:&gt;|>)|!(?:&gt;|>))\s*/);
      if (markerMatch) {
        const markerRaw = markerMatch[1]; // например "~&gt;" или "->"
        const markerType = markerRaw[0];  // "~" | "-" | "!"
        html = html.slice(markerMatch[0].length);
        paragraph.innerHTML = html;

        paragraph.classList.add('note-block');
        if (markerType === '~') {
          noteColor = 'yellow';
          paragraph.classList.add('note-block-yellow');
          iconHTML = '<i class="fas fa-exclamation-triangle" style="color:#ffc107;margin-right:.5em;"></i>';
        } else if (markerType === '-') {
          noteColor = 'blue';
          paragraph.classList.add('note-block-blue');
          iconHTML = '<i class="fas fa-info-circle" style="color:#007acc;margin-right:.5em;"></i>';
        } else if (markerType === '!') {
          noteColor = 'red';
          paragraph.classList.add('note-block-red');
          iconHTML = '<i class="fas fa-exclamation-circle" style="color:#d9534f;margin-right:.5em;"></i>';
        }
      }

      // Первый непустой узел
      const firstNonWsNode = Array.from(paragraph.childNodes).find(function (n) {
        return !(n.nodeType === 3 && /^\s*$/.test(n.nodeValue));
      });

      let movedHTML = '';

      // Переносим только если это самый первый <strong> и он помечен make-heading (или по разделителю, если включено)
      if (firstNonWsNode && firstNonWsNode.nodeType === 1 && firstNonWsNode.tagName === 'STRONG') {
        const strongEl = firstNonWsNode;
        const hasKeep = strongEl.classList.contains('keep-bold') || strongEl.hasAttribute('data-keep-bold');
        const hasMake = strongEl.classList.contains('make-heading') || strongEl.hasAttribute('data-make-heading');
        const endsWithSep = CONFIG.autoBySeparator && CONFIG.separatorRe.test((strongEl.textContent || '').trim());

        if (!hasKeep && (hasMake || endsWithSep)) {
          const clone = strongEl.cloneNode(true);
          if (CONFIG.autoBySeparator) clone.textContent = (clone.textContent || '').replace(CONFIG.separatorRe, '');
          if (noteColor) clone.classList.add('text-' + noteColor);
          movedHTML = clone.outerHTML + '<br><br>';
          strongEl.remove();
        }
      }
      // Покрашиваем оставшиеся <strong> внутри note-блоков (кроме keep-bold)
      if (noteColor) {
        paragraph.querySelectorAll('strong').forEach(function (s) {
          if (s.classList.contains('keep-bold') || s.hasAttribute('data-keep-bold')) return;
          s.classList.add('text-' + noteColor);
        });
      }

      // Вставляем иконку и, если есть, вынесенный заголовок
      if (iconHTML || movedHTML) {
        paragraph.insertAdjacentHTML('afterbegin', iconHTML + movedHTML);
      }

      // Если это note-блок и заголовок не переносили, но в начале стоит <strong> со словом Note/Warning/...
      // — вставим пустую строку после этого <strong>
      if (
        noteColor &&
        !movedHTML &&
        firstNonWsNode &&
        firstNonWsNode.nodeType === 1 &&
        firstNonWsNode.tagName === 'STRONG' &&
        !(firstNonWsNode.classList.contains('keep-bold') || firstNonWsNode.hasAttribute('data-keep-bold')) &&
        CONFIG.addGapAfterKeywordStrongInNotes
      ) {
        const txt = (firstNonWsNode.textContent || '').trim().replace(CONFIG.separatorRe, '');
        const isKeyword = CONFIG.gapKeywords.some(k => k.toLowerCase() === txt.toLowerCase());
        if (isKeyword) {
          // Проверим, нет ли уже <br><br> сразу после
          let n = firstNonWsNode.nextSibling;
          while (n && n.nodeType === 3 && /^\s*$/.test(n.nodeValue)) n = n.nextSibling;
          const firstIsBr = n && n.nodeType === 1 && n.tagName === 'BR';
          const secondIsBr = firstIsBr && n.nextSibling && n.nextSibling.nodeType === 1 && n.nextSibling.tagName === 'BR';

          if (!firstIsBr) {
            firstNonWsNode.insertAdjacentHTML('afterend', '<br><br>');
          } else if (!secondIsBr) {
            n.insertAdjacentHTML('afterend', '<br>');
          }
        }
      }

      paragraph.dataset.boldjsProcessed = '1';
    });
  }

  // instant navigation выключен — запускаем после загрузки
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', processParagraphs);
  } else {
    processParagraphs();
  }
})();
