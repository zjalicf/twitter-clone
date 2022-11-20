import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { MainPageComponent } from './components/main-page/main-page.component';
import { RegisterBusinessComponent } from './components/register-business/register-business.component';
import { RegisterRegularComponent } from './components/register-regular/register-regular.component';

const routes: Routes = [
  {
    path: 'Main-Page',
    component: MainPageComponent
  },
  {
    path: 'Register-Regular',
    component: RegisterRegularComponent
  },
  {
    path: 'Register-Business',
    component: RegisterBusinessComponent
  },
  {
    path: 'Login',
    component: LoginComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
