import React from 'react';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import LoginPage from './components/LoginPage';
import GroupsPage from './components/GroupsPage';
import ContestsPage from './components/ContestsPage';
import ContestDetails from './components/ContestDetails';

const App = () => {
    return (
        <Router>
            <Routes>
                {this.localStorage.getItem("isAuthorized") === false &&
                    <Route path="/" element={<LoginPage/>}/>
                }
                {this.localStorage.getItem('isAuthorized') &&
                    <Route path="/groups" element={<GroupsPage/>}/>
                }
                {this.localStorage.getItem('isAuthorized') &&
                    <Route path="/contests/:groupCode" element={<ContestsPage/>}/>
                }
                {this.localStorage.getItem('isAuthorized') &&
                    <Route path="/contest-details/:groupCode/:contestId" element={<ContestDetails/>}/>
                }

            </Routes>
        </Router>
    )
        ;
};

export default App;
