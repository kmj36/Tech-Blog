import React from 'react';
import './editor.css'
import axios from 'axios';

function Toolbars(set_data) {
    if (set_data === 'h1' || set_data === 'h2' || set_data === 'h3' || set_data === 'p') { 
      document.execCommand('formatBlock', false, set_data);
    } else if(set_data === 'foreColor'){
      var colordata = document.getElementById("RGBcode").value
      document.execCommand(set_data, false, colordata);
    } else if (set_data === 'fontName'){
      var fontvar = document.getElementById("fonts");
      document.execCommand(set_data, false, fontvar.options[fontvar.selectedIndex].value)
    } else if(set_data === 'fontSize'){
      var sizevar = document.getElementById("fontchange");
      document.execCommand(set_data, false, sizevar.options[sizevar.selectedIndex].value);
    } else {
      document.execCommand(set_data); 
    }
}

function Editorhtml() {

  const selectedimage = (event) => {
    var reader = new FileReader();
    reader.onload = function() {
      var result = reader.result;
      var form = new FormData();
      form.append("imagebase64", result);
      axios.post(/*"http://127.0.0.1:36530/api/imageup"*/"https://anend.site:53373/api/imageup", form).then(res => {
        var imagetag = document.createElement("img");
        imagetag.setAttribute("src", res.data.URL);
        document.getElementById("editor").appendChild(imagetag);
      });
    };
    reader.readAsDataURL(event.target.files[0])
  }

  return(
    <div id = "post_editor_container">
    <div className = "editor_ribbon_container">
      <div className = "ercwrapper">
        <div className = "fontedit">
        <select id = "fonts" onInput={() => Toolbars('fontName')}>
            <option value="default">default</option>
            <option value="굴림">굴림</option>
            <option value="굴림체">굴림체</option>
            <option value="궁서">궁서</option>
            <option value="궁서체">궁서체</option>
            <option value="돋움">돋움</option>
            <option value="돋움체">돋움체</option>
            <option value="바탕">바탕</option>
            <option value="바탕체">바탕체</option>
            <option value="휴먼엽서체">휴먼엽서체</option>
            <option value="Nanum Myeongjo">나눔명조</option>
            <option value="Nanum Gothic">나눔고딕</option>
            <option value="Andale Mono">Andale Mono</option>
            <option value="Arial">Arial</option>
            <option value="Arial Black">Arial Black</option>
            <option value="Arial Narrow">Arial Narrow</option>
            <option value="Bookman Old Style">Bookman Old Style</option>
            <option value="Noto Sans">Noto Sans</option>
            <option value="Copperlate Gothic">Copperlate Gothic</option>
            <option value="Courier">Courier</option>
            <option value="Courier New">Courier New</option>
            <option value="Fixedsys">Fixedsys</option>
            <option value="Garamond">Garamond</option>
            <option value="Georgia">Georgia</option>
            <option value="Impact">Impact</option>
            <option value="Lucida Blackletter">Lucida Blackletter</option>
            <option value="Lucida Bright">Lucida Bright</option>
            <option value="Lucida Calligraphy Italic">Lucida Calligraphy Italic</option>
            <option value="Lucida Console">Lucida Console</option>
            <option value="Map Symbols">Map Symbols</option>
            <option value="Marlett">Marlett</option>
            <option value="MingLiu">MingLiu</option>
            <option value="Minion Web">Minion Web</option>
            <option value="Modem">Modem</option>
            <option value="Monotype Sorts">Monotype Sorts</option>
            <option value="Monotype.com">Monotype.com</option>
            <option value="MS Gothic">MS Gothic</option>
            <option value="MS Hei">MS Hei</option>
            <option value="MS Outlook">MS Outlook</option>
            <option value="MS Sans Serif">MS Sans Serif</option>
            <option value="MS Serif">MS Serif</option>
            <option value="MS Song">MS Song</option>
            <option value="MS-DOS CP">MS-DOS CP</option>
            <option value="MT Extra">MT Extra</option>
            <option value="Papyrus">Papyrus</option>
            <option value="Poor Richard">Poor Richard</option>
            <option value="small Fonts">small Fonts</option>
            <option value="Symbol">Symbol</option>
            <option value="System">System</option>
            <option value="Tahoma">Tahoma</option>
            <option value="Terminal">Terminal</option>
            <option value="Times New Roman">Times New Roman</option>
            <option value="Terbuche MS">Terbuche MS</option>
            <option value="Verdana">Verdana</option>
            <option value="Verdata Italic">Verdata Italic</option>
            <option value="Viner Hand ITC">Viner Hand ITC</option>
            <option value="Webdings">Webdings</option>
            <option value="WingDings">WingDings</option>
          </select>
        <select id = "fontchange" onInput={() => Toolbars('fontSize')}>
          <option value="">폰트 크기</option>
          <option value="1">4px</option>
          <option value="2">8px</option>
          <option value="3">10px</option>
          <option value="4">12px</option>
          <option value="5">16px</option>
          <option value="6">20px</option>
          <option value="7">30px</option>
        </select>
          <button className="button" onClick={() => Toolbars('foreColor')}><img className = "palette" src="Palette.png" alt="Palette"/></button>
          <input id = "RGBcode" type="text" placeholder="#000000"></input>
        </div>
        <div className="align_container">
          <button className="button" onClick={() => Toolbars('justifyLeft')}><img className = "alignbutton" src="left.png" alt="left"/></button>
          <button className="button" onClick={() => Toolbars('justifyCenter')}><img className = "alignbutton" src="middle.png" alt="middle"/></button>
          <button className="button" onClick={() => Toolbars('justifyRight')}><img className = "alignbutton" src="right.png" alt="right"/></button>
        </div>
        <div className="imageupload">
          <p className="inline">이미지 넣기</p>
          <img className = "imgcon" src="image_con.png" alt="upload"/>
          <input type="file" name="selimg" accept="image/*" onChange={selectedimage}/>
        </div>
        <div className = "fontdesign">
          <button className="button" onClick={() => Toolbars('bold')}><b>Bold</b></button>
          <button className="button" onClick={() => Toolbars('italic')}><i>Italic</i></button>
          <button className="button" onClick={() => Toolbars('underline')}><u>A</u></button>
          <button className="button" onClick={() => Toolbars('strikeThrough')}><strike>A</strike></button>
        </div>
      </div>
    </div>
    <div className = "editor_body_container">
      <div className = "ebcwrapper">
        <div id = "editor" contentEditable="true"></div>
      </div>
    </div>
  </div>
    );
}

export default Editorhtml;
