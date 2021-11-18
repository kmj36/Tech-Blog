import styled from 'styled-components';
import './smallBox.css';

const NullImage = styled.div`
    display: inline-block;
    width: 60px;
    height: 60px;
    border-style: solid;
    border-width: 1px;
    margin: 5px;
    background-color: grey;
`;

function Searchform(result) {
    return (
        <div className="searchform">
            <a className="openpost" href={"/pages?id=" + result.id}>
                <div className="searchbox">
                    <div className="searchleft">
                        <p className="searchID">{result.id}</p>
                        {result.thumbURL === "" ? <NullImage /> : <img className="searchthumbimgset" src={result.thumbURL} alt='' />}
                    </div>
                    <h3 className="searchtitle">{result.title}</h3>
                    <div className="searchright">
                        <p>조회수: {result.views}</p>
                        <p className="searchdate">{result.postuploaddate}</p>
                    </div>
                </div>
            </a>
        </div>
    )
}

export default Searchform;