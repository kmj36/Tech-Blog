import { useEffect, useState } from 'react';
import SmallBox from '../forms/smallBox';
import axios from 'axios';
import './searchpost.css';

function Searchpost(searchdata){
    const [search, setsearch] = useState();

    const getsearch = async () => {
        const res = await axios.get(/*"http://127.0.0.1:36530/api/posts?name="*/"https://anend.site:53373/api/posts?name=" + searchdata);
        setsearch(res.data);
    };

    useEffect(() => {getsearch();}, []);
    return(
        <div className = "searchwrapper">
            <div className = "searchcounttop">
                <h4 className = "searchcount">검색 결과: {search?.length}</h4>
                </div>
            <div className = "searchbody">
                {search?.map((result) => SmallBox(result))}
            </div>
        </div>
    )
}

export default Searchpost;