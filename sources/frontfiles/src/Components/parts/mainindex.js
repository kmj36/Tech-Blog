import React, { useEffect, useState } from 'react';
import axios from 'axios';
import Postbox from '../forms/postbox';
import './mainindex.css';

function Mainpagecode() {
    const [postlist, setPostlist] = useState();

    const getpostlist = async () => {
        const res = await axios.get(/*"http://127.0.0.1:36530/api/posts"*/"https://anend.site:53373/api/posts");
        setPostlist(res.data);
    };

    useEffect(() => { getpostlist(); }, []);

    return (
        <div className="postwrapper">
            <div className="main"><h4 className="topbarmain">Home</h4></div>
            {postlist?.map((data) => (Postbox(data)))}
        </div>
    );
};

export default Mainpagecode;