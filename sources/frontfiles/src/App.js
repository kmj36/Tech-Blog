import './all.css';
import Topbar from './Components/parts/topbar';
import About from './Pages/About';
import Main from './Pages/Main';
import Search from './Pages/Search';
import Board from './Pages/Board';
import Upload from './Pages/Upload';
import Pages from './Pages/Pages';
import NotFound from './Pages/NotFound'
import { Route, Switch } from 'react-router-dom';


function App() {
  document.title = "Anend BL0G";
  document.oncontextmenu = function(){return false;}
  return(
    <>
    <Switch>
      <Route exact path="/" component={Main}/>
      <Route exact path="/about" component={About}/>
      <Route exact path="/search" component={Search}/>
      <Route exact path="/board" component={Board}/>
      <Route exact path="/upload" component={Upload}/>
      <Route exact path="/pages" component = {Pages}/>
      <Route component={NotFound}/>
    </Switch>
    <Topbar/>
    </>
  );
}

export default App;
