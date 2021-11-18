import React, {useEffect, useState} from 'react';
import axios from 'axios';
import Smallbox from '../forms/smallBox';
import './more.css';

function More() {
    const [count, setcount] = useState(1);
    const [data, setdata] = useState();

    const getdata = async (num) => {
        const jsondata = await axios.get(/*"http://127.0.0.1:36530/api/posts?page="*/"https://anend.site:53373/api/posts?page=" + num)
        setdata(jsondata.data);
    }

    useEffect(()=> {getdata(count);}, [])
    
    let nextbutton, prevbutton;
    if (count === 1 && data?.length < 20 || data?.length === undefined) { // 없음
        prevbutton = {
            visibility: "hidden"
        };
        nextbutton = {
            visibility: "hidden"
        };
    }else if(count === 1) { // 다음
        prevbutton = {
            visibility: "hidden"
        };
    }else if(data?.length < 20) { // 이전
        nextbutton = {
            visibility: "hidden"
        };
    }

    return(
        <div className = "Morepagewrapper">
            <div className = "Moretopbar">
             <h4 className = "Morestyle">more</h4>
            </div>
            <button style={prevbutton} onClick={() => {getdata(count-1); setcount(count-1); }}>이전</button>
            <button style={nextbutton} onClick={() => {getdata(count+1); setcount(count+1); }}>다음</button>
            {data?.map((result) => Smallbox(result))}
            <button style={prevbutton} onClick={() => {getdata(count-1); setcount(count-1); }}>이전</button>
            <button style={nextbutton} onClick={() => {getdata(count+1); setcount(count+1); }}>다음</button>
        </div>
    )
}

export default More;