import Cookies from 'js-cookie'

export default function logout() {
    localStorage.clear();
    sessionStorage.clear();
    Cookies.remove('userKey');
    Cookies.remove('userSecret')
}
