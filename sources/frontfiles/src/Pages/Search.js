import Body from '../Components/parts/body';
import Searchpost from '../Components/parts/searchpost';
import qs from 'query-string';

function Search(location) {
    const searchdata = qs.parse(location.location.search).search_post
    return(
        <>
        {Body(Searchpost(searchdata))}
        </>
    );
}

export default Search;