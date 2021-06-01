import React from 'react';
import './topbar.css';

function Topbar_main() { // 컨테이너는 위치가 고정이니 postion으로 설정하는 방법도 있음
    return(
        <div className = "header">
            <div className = "hdwrapper">
                    <div className = "explain_container">
                        <a href="/" className="astyle">
                            <img src = "K.svg" alt = "logo" className = "logo" />
                            <h1 className = "banner">Anend</h1>
                        </a>
                    </div>
                    <div className = "search_container">
                        <form className = "search_form" name = "postsearch" action = "/search" acceptCharset = "utf-8" method = "GET">
                            <input className = "textbox" type = "text" name = "search_post" alt = "search" minLength="2" required></input>
                            <input className = "buttonbox" type = "submit" value = "검색" alt = "searchsummbit" />
                        </form>
                    </div>
                    <div className = "navbar_container">
                        <div className = "rightset">
                            <div className = "nav_b"><a href = "/" className = "astyle">홈</a></div>
                            <div className = "nav_b"><a href = "/board" className = "astyle">더보기</a></div>
                            <div className = "nav_b"><a href = "/about" className = "astyle">소개</a></div>
                        </div>
                    </div>
                </div>
        </div>
    );
}

export default Topbar_main;