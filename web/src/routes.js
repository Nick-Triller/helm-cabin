import Releases from './components/ReleasesComponent.vue'
import NotFoundComponent from "./components/NotFoundComponent";
import ReleaseComponent from "./components/ReleaseComponent";

const routes = [
  { path: '/', component: Releases },
  { path: '/releases/:name', component: ReleaseComponent},
  { path: '*', component: NotFoundComponent }
];

export default routes
