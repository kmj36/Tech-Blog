import React from 'react';
import styled from 'styled-components';
import './postbox.css';

const NullImage = styled.div`
    display: inline-block;
    width: 140px;
    height: 140px;
    border-style: solid;
    border-width: 1px;
    margin: 5px;
    background-color: grey;
`;

function Postbox(jsondata) {
    if (jsondata.secret === false) {
        return (
            <div className="postform">
                <a className="openpost" href={"/pages?id=" + jsondata.id}>
                    {jsondata.thumbURL === "" ? <NullImage /> : <img className="thumbimgset" src={jsondata.thumbURL} alt='' />}
                    <div className="postinfobox">
                        <h2 className="posttitle">{jsondata.title}</h2>
                        <p className="postbody">{jsondata.body}...</p>
                        <p className="postdate">{jsondata.postuploaddate}</p>
                    </div>
                </a>
            </div>
        );
    } else {
        return (
            <div className="postform">
                <a className="openpost" href={"/pages?id=" + jsondata.id}>
                    {jsondata.thumbURL === "" ? <NullImage /> : <img className="thumbimgset" src={jsondata.thumbURL} alt='' />}
                    <div className="postinfobox">
                        <h2 className="posttitle">{jsondata.title}</h2>
                        <p className="postbody">...</p>
                        <p className="postdate">{jsondata.postuploaddate}</p>
                        <p className="postbemil">비밀글</p>
                    </div>
                </a>
            </div>
        );
    }
}

export default Postbox;